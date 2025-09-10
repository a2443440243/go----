package service

import (
	"go-api-scaffold/internal/repository"
	"go-api-scaffold/pkg/logger"
)

// Service 服务接口集合
type Service struct {
	User UserService
}

// New 创建服务实例
func New(repo *repository.Repository, logger logger.Logger) *Service {
	return &Service{
		User: NewUserService(repo.User, logger),
	}
}
