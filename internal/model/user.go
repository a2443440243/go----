package model

import (
	"time"
)

// User 用户模型
type User struct {
	BaseModel
	Username  string     `json:"username" gorm:"uniqueIndex;not null;size:50" validate:"required,min=3,max=50"`
	Email     string     `json:"email" gorm:"uniqueIndex;not null;size:100" validate:"required,email"`
	Password  string     `json:"-" gorm:"not null;size:255" validate:"required,min=6"`
	Nickname  string     `json:"nickname" gorm:"size:50"`
	Avatar    string     `json:"avatar" gorm:"size:255"`
	Phone     string     `json:"phone" gorm:"size:20"`
	Status    int        `json:"status" gorm:"default:1;comment:用户状态 1:正常 0:禁用"`
	LastLogin *time.Time `json:"last_login"`
}

// TableName 表名
func (User) TableName() string {
	return "users"
}

// UserCreateRequest 创建用户请求
type UserCreateRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Nickname string `json:"nickname" validate:"max=50"`
	Phone    string `json:"phone" validate:"max=20"`
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	Nickname string `json:"nickname" validate:"max=50"`
	Avatar   string `json:"avatar" validate:"max=255"`
	Phone    string `json:"phone" validate:"max=20"`
	Status   *int   `json:"status" validate:"omitempty,oneof=0 1"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Phone     string    `json:"phone"`
	Status    int       `json:"status"`
	LastLogin time.Time `json:"last_login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
