package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"go-api-scaffold/internal/model"
	"go-api-scaffold/internal/service"
	"go-api-scaffold/pkg/logger"
	"go-api-scaffold/pkg/response"
)

// UserHandler 用户处理器
type UserHandler struct {
	service   service.UserService
	logger    logger.Logger
	validator *validator.Validate
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(service service.UserService, logger logger.Logger) *UserHandler {
	return &UserHandler{
		service:   service,
		logger:    logger,
		validator: validator.New(),
	}
}

// Create 创建用户
func (h *UserHandler) Create(c *gin.Context) {
	var req model.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body: " + err.Error())
		response.BadRequest(c, "Invalid request body")
		return
	}

	// 验证请求参数
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Validation failed: " + err.Error())
		response.BadRequest(c, "Validation failed: "+err.Error())
		return
	}

	// 创建用户
	user, err := h.service.Create(&req)
	if err != nil {
		h.logger.Error("Failed to create user: " + err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "User created successfully", user)
}

// GetByID 根据ID获取用户
func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	user, err := h.service.GetByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get user: " + err.Error())
		if err.Error() == "user not found" {
			response.NotFound(c, "User not found")
		} else {
			response.ServerError(c, "Failed to get user")
		}
		return
	}

	response.Success(c, user)
}

// Update 更新用户
func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req model.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body: " + err.Error())
		response.BadRequest(c, "Invalid request body")
		return
	}

	// 验证请求参数
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Validation failed: " + err.Error())
		response.BadRequest(c, "Validation failed: "+err.Error())
		return
	}

	// 更新用户
	user, err := h.service.Update(uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update user: " + err.Error())
		if err.Error() == "user not found" {
			response.NotFound(c, "User not found")
		} else {
			response.BadRequest(c, err.Error())
		}
		return
	}

	response.SuccessWithMessage(c, "User updated successfully", user)
}

// Delete 删除用户
func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		h.logger.Error("Failed to delete user: " + err.Error())
		if err.Error() == "user not found" {
			response.NotFound(c, "User not found")
		} else {
			response.ServerError(c, "Failed to delete user")
		}
		return
	}

	response.SuccessWithMessage(c, "User deleted successfully", nil)
}

// List 获取用户列表
func (h *UserHandler) List(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 限制分页大小
	if pageSize > 100 {
		pageSize = 100
	}

	users, meta, err := h.service.List(page, pageSize)
	if err != nil {
		h.logger.Error("Failed to get user list: " + err.Error())
		response.ServerError(c, "Failed to get user list")
		return
	}

	response.SuccessPage(c, users, *meta)
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body: " + err.Error())
		response.BadRequest(c, "Invalid request body")
		return
	}

	// 验证请求参数
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Validation failed: " + err.Error())
		response.BadRequest(c, "Username and password are required")
		return
	}

	// 用户登录
	user, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		h.logger.Error("Login failed: " + err.Error())
		response.Unauthorized(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Login successful", user)
}
