package start

import (
	"context"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/mr-linch/go-tg/tgb/session"
	"tacy/internal/models"
	"tacy/internal/services/authService"
)

type StartHandler struct {
	sessionManager *session.Manager[models.Session]

	AuthService authService.Service
}

func NewStartHandler(sm *session.Manager[models.Session], s authService.Service) *StartHandler {
	return &StartHandler{
		sessionManager: sm, AuthService: s,
	}
}

func (s *StartHandler) Start(ctx context.Context, msg *tgb.MessageUpdate) error {
	role, err := s.AuthService.GetRoleById(ctx, msg.Update.Message.From.ID)
	if err != nil || role == 0 {
		return msg.Answer("Error").DoVoid(ctx)
	}

	switch role {
	case models.AdminRole:
		s.sessionManager.Get(ctx).Step = models.SessionStepAdminMenuHandler

		return msg.Answer("Пожалуйста, выберите действие").
			ReplyMarkup(buildAdminStartMenu()).
			DoVoid(ctx)
	case models.UserRole:
		s.sessionManager.Get(ctx).Step = models.SessionStepUserMenuHandler

		return msg.Answer("Пожалуйста, выберите действие").
			ReplyMarkup(buildUserStartMenu()).
			DoVoid(ctx)
	default:
		return nil
	}
}
