package admin

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/mr-linch/go-tg/tgb/session"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"tacy/internal/models"
	"tacy/internal/services/userService"
)

type AdminHandler struct {
	sessionManager *session.Manager[models.Session]
	UserService    userService.Service
	logger         zerolog.Logger
	Client         *tg.Client
}

func NewAdminHandler(
	sm *session.Manager[models.Session],
	UserService userService.Service,
	logger zerolog.Logger,
	client *tg.Client,
) *AdminHandler {
	return &AdminHandler{
		sessionManager: sm,
		UserService:    UserService,
		logger:         logger,
		Client:         client,
	}
}

func (h *AdminHandler) AdminStartMenuHandler(ctx context.Context, msg *tgb.MessageUpdate) error {
	switch msg.Text {
	case models.AdminStartMenu.AdminEnter:
		h.logger.Info().Msg(msg.Text + " fetched.")
		h.sessionManager.Get(ctx).Step = models.SessionStepEnterAdminMenuHandler
		kb := tg.NewReplyKeyboardMarkup(
			tg.NewButtonColumn(
				tg.NewKeyboardButton(models.AdminMenu.InsertNewPhoto),
				tg.NewKeyboardButton(models.AdminMenu.InsertNewCompliment),
				tg.NewKeyboardButton(models.AdminMenu.GetAllPhotos),
				tg.NewKeyboardButton(models.AdminMenu.GetAllCompliments),
			)...,
		).WithResizeKeyboardMarkup()
		return msg.Update.Reply(ctx, msg.Answer(fmt.Sprintf("Выберите действие")).
			ReplyMarkup(kb))
	default:
		h.logger.Warn().Msg(msg.Text + " fetched. Unknown endpoint.")
		h.sessionManager.Get(ctx).Step = models.SessionStepInit
		return msg.Answer("Напишите /start").ReplyMarkup(tg.NewReplyKeyboardRemove()).DoVoid(ctx)
	}
}

func (h *AdminHandler) AdminMenuHandler(ctx context.Context, msg *tgb.MessageUpdate) error {
	switch msg.Text {
	case models.AdminMenu.InsertNewPhoto:
		h.sessionManager.Get(ctx).Step = models.SessionStepInsertNewPhotoMenuHandler
		kb := tg.NewReplyKeyboardMarkup(
			tg.NewButtonColumn(
				tg.NewKeyboardButton(models.AdminStartMenu.AdminEnter),
			)...,
		).WithResizeKeyboardMarkup()
		return msg.Update.Reply(ctx, msg.Answer(fmt.Sprintf("Пришлите фото, которое вы ходите добавить")).
			ReplyMarkup(kb))
	case models.AdminMenu.InsertNewCompliment:
		h.sessionManager.Get(ctx).Step = models.SessionStepInsertNewComplimentMenuHandler
		kb := tg.NewReplyKeyboardMarkup(
			tg.NewButtonColumn(
				tg.NewKeyboardButton(models.AdminStartMenu.AdminEnter),
			)...,
		).WithResizeKeyboardMarkup()
		return msg.Update.Reply(ctx, msg.Answer(fmt.Sprintf("Пришлите комплимент, который вы хотите добавить")).
			ReplyMarkup(kb))

	case models.AdminMenu.GetAllCompliments:
		return h.AdminMenuGetAllPhotosHandler(ctx, msg)

	case models.AdminMenu.GetAllPhotos:
		return h.AdminMenuGetAllPhotosHandler(ctx, msg)
	default:
		h.logger.Info().Msg(msg.Text + " fetched. Unknown endpoint")
		h.sessionManager.Get(ctx).Step = models.SessionStepInit
		return msg.Answer("Напишите /start").ReplyMarkup(tg.NewReplyKeyboardRemove()).DoVoid(ctx)
	}
}

func (h *AdminHandler) AdminMenuInsertPhotoHandler(ctx context.Context, msg *tgb.MessageUpdate) error {
	if msg.Photo != nil {
		h.logger.Info().Msg("Photo" + " fetched.")
		photo := msg.Message.Photo[len(msg.Message.Photo)-1]
		file, err := h.Client.GetFile(photo.FileID).Do(ctx)
		if err != nil {
			h.logger.Warn().Err(err)
			return msg.Update.Reply(ctx, msg.Answer("unable"))
		}
		if file.FileSize == 0 {
			h.logger.Warn().Msg("file is nil")
			return msg.Update.Reply(ctx, msg.Answer("unable"))
		}
		photoRC, err := h.Client.Download(ctx, file.FilePath)
		if err != nil {
			h.logger.Warn().Err(err)
			return msg.Update.Reply(ctx, msg.Answer("Unable to add photo."))
		}
		defer func() {
			if err := photoRC.Close(); err != nil {
				h.logger.Warn().Err(err)
			}
		}()
		bytes, err := io.ReadAll(photoRC)
		if err != nil {
			h.logger.Warn().Err(err)
			return msg.Update.Reply(ctx, msg.Answer("Извините, фото не было добавлено, попробуйте снова."))
		}
		err = h.UserService.InputPhotoFromAdmin(ctx, bytes)
		if err != nil {
			log.Warn().Err(err).Msg("Photo not added to db.")
			return msg.Update.Reply(ctx, msg.Answer("Извините, фото не было добавлено, попробуйте снова"))
		}
		log.Info().Msg("Photo added.")
		return msg.Update.Reply(ctx, msg.Answer("Фото успешно добавлено."))
	}
	switch msg.Text {
	case models.AdminStartMenu.AdminEnter:
		h.sessionManager.Get(ctx).Step = models.SessionStepEnterAdminMenuHandler
		return h.AdminStartMenuHandler(ctx, msg)
	default:
		h.sessionManager.Get(ctx).Step = models.SessionStepInit
		return msg.Answer("Напишите /start").ReplyMarkup(tg.NewReplyKeyboardRemove()).DoVoid(ctx)
	}
}

func (h *AdminHandler) AdminMenuInsertNewComplimentHandler(ctx context.Context, msg *tgb.MessageUpdate) error {
	//в будущем - добавить update запрос на комплименты
	switch msg.Text {
	case models.AdminStartMenu.AdminEnter:
		h.sessionManager.Get(ctx).Step = models.SessionStepEnterAdminMenuHandler
		return h.AdminStartMenuHandler(ctx, msg)
	default:
		h.sessionManager.Get(ctx).Step = models.SessionStepInit
		return msg.Answer("Напишите /start").ReplyMarkup(tg.NewReplyKeyboardRemove()).DoVoid(ctx)
	}
}

func (h *AdminHandler) AdminMenuGetAllPhotosHandler(ctx context.Context, msg *tgb.MessageUpdate) error {
	photos, err := h.UserService.OutputAllPhotos(ctx)
	if err != nil {
		h.logger.Warn().Err(err).Msg("Unable to get photos.")
		return msg.Update.Reply(ctx, msg.Answer("Unable to get photos."))
	}
	for _, photo := range photos {
		file := tg.NewInputFileBytes("photo", photo)
		msg.Update.Reply(ctx, msg.AnswerPhoto(tg.NewFileArgUpload(file)))
	}
	if err != nil {
		h.logger.Warn().Err(err).Msg("Unable to send photo.")
		return msg.Update.Reply(ctx, msg.Answer("Unable to send photo."))
	}
	if err == nil {
		h.sessionManager.Get(ctx).Step = models.SessionStepEnterAdminMenuHandler
		return msg.Update.Reply(ctx, msg.Answer("Done!"))
	}
	switch msg.Text {
	case models.AdminStartMenu.AdminEnter:
		h.sessionManager.Get(ctx).Step = models.SessionStepEnterAdminMenuHandler
		return h.AdminStartMenuHandler(ctx, msg)
	default:
		h.sessionManager.Get(ctx).Step = models.SessionStepInit
		return msg.Answer("Напишите /start").ReplyMarkup(tg.NewReplyKeyboardRemove()).DoVoid(ctx)
	}
}

func (h *AdminHandler) AdminMenuGetAllComplimentsHandler(ctx context.Context, msg *tgb.MessageUpdate) error {
	switch msg.Text {
	case models.AdminStartMenu.AdminEnter:
		h.sessionManager.Get(ctx).Step = models.SessionStepEnterAdminMenuHandler
		return h.AdminStartMenuHandler(ctx, msg)
	default:
		h.sessionManager.Get(ctx).Step = models.SessionStepInit
		return msg.Answer("Напишите /start").ReplyMarkup(tg.NewReplyKeyboardRemove()).DoVoid(ctx)
	}
}
