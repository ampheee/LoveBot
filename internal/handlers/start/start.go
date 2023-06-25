package start

import (
	"context"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/mr-linch/go-tg/tgb/session"
	"github.com/rs/zerolog"
	"strconv"
	"tacy/internal/models"
	"tacy/internal/services/authService"
	"tacy/pkg/botlogger"
)

type StartHandler struct {
	sessionManager *session.Manager[models.Session]
	logger         zerolog.Logger
	AuthService    authService.Service
}

func NewStartHandler(sm *session.Manager[models.Session], s authService.Service) *StartHandler {
	return &StartHandler{
		sessionManager: sm, AuthService: s, logger: botlogger.GetLogger(),
	}
}

func (s *StartHandler) Start(ctx context.Context, msg *tgb.MessageUpdate) error {
	role, err := s.AuthService.GetRoleById(ctx, msg.Update.Message.From.ID)
	if err != nil || role == 0 {
		s.logger.Warn().Err(err)
		s.logger.Info().Msg("role is " + strconv.Itoa(role))
		return msg.Answer("Непредвиденная ошибка на стороне хоста.").DoVoid(ctx)
	}
	switch role {
	case models.AdminRole:
		s.sessionManager.Get(ctx).Step = models.SessionStepAdminMenuHandler
		return msg.Answer("Пожалуйста, выберите действие").
			ReplyMarkup(buildAdminStartMenu()).
			DoVoid(ctx)
	case models.UserRole:
		s.sessionManager.Get(ctx).Step = models.SessionStepUserMenuHandler
		return msg.Answer("Привет, Тась :)\nИногда я не успеваю напоминать тебе о том, " +
			"насколько ты молодец и подмечать то, как ты стараешься 😢\nИ дабы исправить это, я решил написать вот такого" +
			" простенького телеграм-бота, который каждые 6-часов будет отправлять тебе\nнебольшие факты и слова поддержки ❤️‍🩹" +
			"Если что-то сломается или захочешь изменить таймаут между сообщениями, то дай знать @ampheee," +
			" постараюсь починить его в скорейшие сроки🛠\nНу а теперь - скорее выбирай действие, потыкай все кнопочки!! 🔆").
			ReplyMarkup(buildUserStartMenu()).
			DoVoid(ctx)
	default:
		return msg.Answer("Вы кто такие ? Мы вас не знаем, до свидания 😒").DoVoid(ctx)
	}
}
