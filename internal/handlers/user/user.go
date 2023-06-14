package user

import (
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/mr-linch/go-tg/tgb/session"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"tacy/internal/models"
	"tacy/internal/services/userService"
)

type UserHandler struct {
	sessionManager *session.Manager[models.Session]
	UserService    userService.USUsecase
	Logger         zerolog.Logger
	Client         *tg.Client
}

func NewUserHandler(
	sm *session.Manager[models.Session],
	UserService userService.USUsecase,
	logger zerolog.Logger,
	client *tg.Client,
) *UserHandler {
	return &UserHandler{
		sessionManager: sm,
		UserService:    UserService,
		Logger:         logger,
		Client:         client,
	}
}

func (h *UserHandler) UserStartMenuSelectionHandler(ctx context.Context, msg *tgb.MessageUpdate) error {
	h.Logger.Info().Msg(msg.Text + " fetched.")
	switch msg.Text {
	case models.UserStartMenu.GetComplimentNow:
		return h.UserGetComplimentByRandomHandler(ctx, msg)
	case models.UserStartMenu.InsertSomeThoughts:
		h.sessionManager.Get(ctx).Step = models.SessionStepInsertSomeThoughts
		return msg.Answer("Оставь тут все, что хочешь, я обязательно это прочту \U0001F979💛").DoVoid(ctx)
	default:
		h.sessionManager.Get(ctx).Step = models.SessionStepInit
		return msg.Answer("Что-то сломалось 😢\nНапиши /start").ReplyMarkup(tg.NewReplyKeyboardRemove()).DoVoid(ctx)
	}
}

func (h *UserHandler) UserGetComplimentByRandomHandler(ctx context.Context, msg *tgb.MessageUpdate) error {
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
	return msg.Answer(compliment).DoVoid(ctx)
}

func (h *UserHandler) UserInputSomeThoughts(ctx context.Context, msg *tgb.MessageUpdate) error {
	if msg.Photo != nil {
		h.Logger.Info().Msg("Photo from user" + " fetched.")
		photo := msg.Message.Photo[len(msg.Message.Photo)-1]
		file, err := h.Client.GetFile(photo.FileID).Do(ctx)
		if err != nil {
			h.Logger.Warn().Err(err)
		}
		if file.FileSize == 0 {
			h.Logger.Warn().Msg("file is nil")
		}
		photoRC, err := h.Client.Download(ctx, file.FilePath)
		if err != nil {
			h.Logger.Warn().Err(err)
		}
		defer func() {
			if err := photoRC.Close(); err != nil {
				h.Logger.Warn().Err(err)
			}
		}()
		bytes, err := io.ReadAll(photoRC)
		if err != nil {
			h.Logger.Warn().Err(err)
		}
		err = h.UserService.InsertPhoto(ctx, bytes, true)
		if err != nil {
			log.Warn().Err(err).Msg("Photo from user not added to db.")
		}
		h.Logger.Info().Msg("Photo added from user.")
	}
	switch msg.Text {
	default:
		var err error
		h.sessionManager.Get(ctx).Step = models.SessionStepUserMenuHandler
		if msg.Text != "" {
			err = h.UserService.InputThoughtsFromUser(ctx, msg.Text)
		}
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
