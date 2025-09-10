package repository

import (
	"gorm.io/gorm"
)

// Repository 仓储接口集合
type Repository struct {
	User UserRepository
}

// New 创建仓储实例
func New(db *gorm.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}