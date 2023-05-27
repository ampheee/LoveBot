package user

import (
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/mr-linch/go-tg/tgb/session"
	"github.com/rs/zerolog"
	"tacy/internal/models"
	"tacy/internal/services/userService"
	"tacy/pkg/botlogger"
)

type UserHandler struct {
	sessionManager *session.Manager[models.Session]
	UserService    userService.Service
	Logger         zerolog.Logger
}

func NewUserHandler(
	sm *session.Manager[models.Session],
	UserService userService.Service,
) *UserHandler {
	return &UserHandler{
		sessionManager: sm,
		UserService:    UserService,
		Logger:         botlogger.GetLogger(),
	}
}

func (h *UserHandler) UserStartMenuSelectionHandler(ctx context.Context, msg *tgb.MessageUpdate) error {
	h.Logger.Info().Msg(msg.Text + " fetched.")
	switch msg.Text {
	case models.UserStartMenu.GetComplimentNow:
		return h.UserGetComplimentByRandomHandler(ctx, msg)
	case models.UserStartMenu.InsertSomeThoughts:
		h.sessionManager.Get(ctx).Step = models.SessionStepInsertSomeThoughts
		return msg.Answer("Оставь тут все, что хочешь, я обязательно это прочту \U0001F979💛").ReplyMarkup(
			tg.NewReplyKeyboardMarkup(
				tg.NewButtonColumn(
					tg.NewKeyboardButton("Вернуться назад ❤️‍🩹"),
				)...,
			).WithResizeKeyboardMarkup()).DoVoid(ctx)
	default:
		h.sessionManager.Get(ctx).Step = models.SessionStepInit
		return msg.Answer("Что-то сломалось 😢\nНапиши /start").ReplyMarkup(tg.NewReplyKeyboardRemove()).DoVoid(ctx)
	}
}

func (h *UserHandler) UserGetComplimentByRandomHandler(ctx context.Context, msg *tgb.MessageUpdate) error {
	h.Logger.Info().Msg(msg.Text + " fetched.")
	photo, compliment, err := h.UserService.OutputComplimentAndPhotoByRandom(ctx)
	if err != nil {
		h.Logger.Warn().Err(err)
		return msg.Answer("Что-то пошло не так 😡 Напиши мне и я все подправлю 👌🏻").DoVoid(ctx)
	}
	err = msg.AnswerPhoto(tg.NewFileArgUpload(tg.NewInputFileBytes("photo", photo))).DoVoid(ctx)
	if err != nil {
		h.Logger.Warn().Err(err)
		return msg.Answer("Не получается отправить фото 😢").DoVoid(ctx)
	}
	err = msg.Answer(compliment).DoVoid(ctx)
	if err != nil {
		h.Logger.Warn().Err(err)
		return msg.Answer("Не получается отправить комплимент 😢 Но мы все знаем, что даже так ты прекрасна " +
			"\U0001F979").DoVoid(ctx)
	}
	return msg.Answer("Я надеюсь ты рада!!)").DoVoid(ctx)
}

func (h *UserHandler) UserInputSomeThoughts(ctx context.Context, msg *tgb.MessageUpdate) error {
	h.Logger.Info().Msg(msg.Text + " fetched.")
	switch msg.Text {
	case models.UserStartMenu.Back:
		h.sessionManager.Get(ctx).Step = models.SessionStepUserMenuHandler
		return h.UserStartMenuSelectionHandler(ctx, msg)
	default:
		err := h.UserService.InputThoughtsFromUser(ctx, msg.Text)
		if err != nil {
			return msg.Answer("Что-то пошло не так @\nЯ решу эту проблему, " +
				"но не забывай, что личные сообщения никто не отменял").DoVoid(ctx)
		}
		return msg.Answer("Спасибо, солнышко, это важно для меня 💛").ReplyMarkup(tg.NewReplyKeyboardMarkup(
			tg.NewButtonColumn(
				tg.NewKeyboardButton(models.UserStartMenu.GetComplimentNow),
				tg.NewKeyboardButton(models.UserStartMenu.InsertSomeThoughts),
			)...,
		).WithResizeKeyboardMarkup()).DoVoid(ctx)
	}
}
