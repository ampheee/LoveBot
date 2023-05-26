package repository

import (
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/rs/zerolog"
	"tacy/internal/services/authService"
	"tacy/internal/storage"
	"tacy/pkg/botlogger"
)

type AuthService struct {
	logger  zerolog.Logger
	storage storage.Storage
}

func (s AuthService) GetRoleById(ctx context.Context, id tg.UserID) (int, error) {
	role, err := s.storage.GetRoleById(ctx, id)
	if err != nil {
		return 2, err
	}
	return role, nil
}

func InitNewAuthService(s storage.Storage) authService.Service {
	return &AuthService{
		logger:  botlogger.GetLogger(),
		storage: s,
	}
}
