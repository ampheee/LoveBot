package handlers

import (
	"github.com/mr-linch/go-tg/tgb"
	"github.com/mr-linch/go-tg/tgb/session"
	"tacy/internal/models"
)

type SessionManager struct {
	*session.Manager[models.Session]
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		session.NewManager(models.Session{
			Step: models.SessionStepInit,
		}),
	}
}

func (h *Handler) registerSession() {
	h.Router.Use(h.sessionManager)
}

func (h *Handler) isSessionStep(state models.SessionStep) tgb.Filter {
	return h.sessionManager.Filter(func(session *models.Session) bool {
		return session.Step == state
	})
}
