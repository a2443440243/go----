package service

import (
	"errors"
	"fmt"
	"math"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"go-api-scaffold/internal/model"
	"go-api-scaffold/internal/repository"
	"go-api-scaffold/pkg/logger"
	"go-api-scaffold/pkg/response"
)

// UserService 用户服务接口
type UserService interface {
	Create(req *model.UserCreateRequest) (*model.UserResponse, error)
	GetByID(id uint) (*model.UserResponse, error)
	Update(id uint, req *model.UserUpdateRequest) (*model.UserResponse, error)
	Delete(id uint) error
	List(page, pageSize int) ([]*model.UserResponse, *response.PageMeta, error)
	Login(username, password string) (*model.UserResponse, error)
}

// userService 用户服务实现
type userService struct {
	repo   repository.UserRepository
	logger logger.Logger
}

// NewUserService 创建用户服务实例
func NewUserService(repo repository.UserRepository, logger logger.Logger) UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

// Create 创建用户
func (s *userService) Create(req *model.UserCreateRequest) (*model.UserResponse, error) {
	// 检查用户名是否存在
	if _, err := s.repo.GetByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否存在
	if _, err := s.repo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to hash password: %v", err))
		return nil, errors.New("failed to create user")
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Status:   1,
	}

	if err := s.repo.Create(user); err != nil {
		s.logger.Error(fmt.Sprintf("Failed to create user: %v", err))
		return nil, errors.New("failed to create user")
	}

	return s.toUserResponse(user), nil
}

// GetByID 根据ID获取用户
func (s *userService) GetByID(id uint) (*model.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		s.logger.Error(fmt.Sprintf("Failed to get user by ID: %v", err))
		return nil, errors.New("failed to get user")
	}

	return s.toUserResponse(user), nil
}

// Update 更新用户
func (s *userService) Update(id uint, req *model.UserUpdateRequest) (*model.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		s.logger.Error(fmt.Sprintf("Failed to get user by ID: %v", err))
		return nil, errors.New("failed to update user")
	}

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := s.repo.Update(user); err != nil {
		s.logger.Error(fmt.Sprintf("Failed to update user: %v", err))
		return nil, errors.New("failed to update user")
	}

	return s.toUserResponse(user), nil
}

// Delete 删除用户
func (s *userService) Delete(id uint) error {
	if _, err := s.repo.GetByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		s.logger.Error(fmt.Sprintf("Failed to get user by ID: %v", err))
		return errors.New("failed to delete user")
	}

	if err := s.repo.Delete(id); err != nil {
		s.logger.Error(fmt.Sprintf("Failed to delete user: %v", err))
		return errors.New("failed to delete user")
	}

	return nil
}

// List 获取用户列表
func (s *userService) List(page, pageSize int) ([]*model.UserResponse, *response.PageMeta, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	users, total, err := s.repo.List(offset, pageSize)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to get user list: %v", err))
		return nil, nil, errors.New("failed to get user list")
	}

	// 转换为响应格式
	userResponses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = s.toUserResponse(user)
	}

	// 计算分页信息
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	meta := &response.PageMeta{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	return userResponses, meta, nil
}

// Login 用户登录
func (s *userService) Login(username, password string) (*model.UserResponse, error) {
	// 根据用户名或邮箱查找用户
	var user *model.User
	var err error

	// 先尝试用户名
	user, err = s.repo.GetByUsername(username)
	if err != nil {
		// 再尝试邮箱
		user, err = s.repo.GetByEmail(username)
		if err != nil {
			return nil, errors.New("invalid username or password")
		}
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("user account is disabled")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLogin = &now
	if err := s.repo.Update(user); err != nil {
		s.logger.Warn(fmt.Sprintf("Failed to update last login time: %v", err))
	}

	return s.toUserResponse(user), nil
}

// toUserResponse 转换为用户响应格式
func (s *userService) toUserResponse(user *model.User) *model.UserResponse {
	resp := &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Phone:     user.Phone,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.LastLogin != nil {
		resp.LastLogin = *user.LastLogin
	}

	return resp
}