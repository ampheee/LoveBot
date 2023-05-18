package admin

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/mr-linch/go-tg/tgb/session"
	"tacy/internal/models"
	"tacy/internal/services/userService"
)

type AdminHandler struct {
	sessionManager *session.Manager[models.Session]
	UserService    userService.Service
}

func NewAdminHandler(
	sm *session.Manager[models.Session],
	UserService userService.Service,
) *AdminHandler {
	return &AdminHandler{
		sessionManager: sm,
		UserService:    UserService,
	}
}

func (h *AdminHandler) AdminStartMenuHandler(ctx context.Context, msg *tgb.MessageUpdate) error {
	switch msg.Text {
	case models.AdminStartMenu.AdminEnter:
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
				tg.NewKeyboardButton(models.AdminMenu.InsertNewPhoto),
				tg.NewKeyboardButton(models.AdminStartMenu.AdminEnter),
			)...,
		).WithResizeKeyboardMarkup()
		return msg.Update.Reply(ctx, msg.Answer(fmt.Sprintf("Пришлите фото, которое вы ходите добавить")).
			ReplyMarkup(kb))
	case models.AdminMenu.InsertNewCompliment:
		h.sessionManager.Get(ctx).Step = models.SessionStepInsertNewComplimentMenuHandler
		kb := tg.NewReplyKeyboardMarkup(
			tg.NewButtonColumn(
				tg.NewKeyboardButton(models.AdminMenu.InsertNewCompliment),
				tg.NewKeyboardButton(models.AdminStartMenu.AdminEnter),
			)...,
		).WithResizeKeyboardMarkup()
		return msg.Update.Reply(ctx, msg.Answer(fmt.Sprintf("Пришлите комплимент, который вы хотите добавить")).
			ReplyMarkup(kb))

	case models.AdminMenu.GetAllCompliments:
		h.sessionManager.Get(ctx).Step = models.SessionStepGetAllComplimentsHandler
		kb := tg.NewReplyKeyboardMarkup(
			tg.NewButtonColumn(
				tg.NewKeyboardButton(models.AdminStartMenu.AdminEnter),
			)...,
		).WithResizeKeyboardMarkup()
		return msg.Update.Reply(ctx, msg.Answer(fmt.Sprintf("Ожидайте")).
			ReplyMarkup(kb))

	case models.AdminMenu.GetAllPhotos:
		h.sessionManager.Get(ctx).Step = models.SessionStepGetAllPhotosMenuHandler
		kb := tg.NewReplyKeyboardMarkup(
			tg.NewButtonColumn(
				tg.NewKeyboardButton(models.AdminStartMenu.AdminEnter),
			)...,
		).WithResizeKeyboardMarkup()
		return msg.Update.Reply(ctx, msg.Answer(fmt.Sprintf("Ожидайте")).
			ReplyMarkup(kb))
	default:
		h.sessionManager.Get(ctx).Step = models.SessionStepInit
		return msg.Answer("Напишите /start").ReplyMarkup(tg.NewReplyKeyboardRemove()).DoVoid(ctx)
	}
}

func (h *AdminHandler) AdminMenuInsertPhotoHandler(ctx context.Context, msg *tgb.MessageUpdate) error {
	if msg.Photo != nil {

	}
	switch msg.Text {
	case models.AdminMenu.InsertNewPhoto:
		h.sessionManager.Get(ctx).Step = models.SessionStepInsertNewPhotoMenuHandler
		kb := tg.NewReplyKeyboardMarkup(
			tg.NewButtonColumn(
				tg.NewKeyboardButton(models.AdminMenu.InsertNewPhoto),
				tg.NewKeyboardButton(models.AdminStartMenu.AdminEnter),
			)...,
		).WithResizeKeyboardMarkup()
		return msg.Update.Reply(ctx, msg.Answer(fmt.Sprintf("Пришлите фото, которое вы ходите добавить")).
			ReplyMarkup(kb))
	case models.AdminStartMenu.AdminEnter:
		h.sessionManager.Get(ctx).Step = models.SessionStepEnterAdminMenuHandler
		return h.AdminStartMenuHandler(ctx, msg)
	default:
		h.sessionManager.Get(ctx).Step = models.SessionStepInit
		return msg.Answer("Напишите /start").ReplyMarkup(tg.NewReplyKeyboardRemove()).DoVoid(ctx)
	}
}

func (h *AdminHandler) AdminMenuInsertNewComplimentHandler(ctx context.Context, msg *tgb.MessageUpdate) error {
	switch msg.Text {
	case models.AdminMenu.InsertNewCompliment:
		h.sessionManager.Get(ctx).Step = models.SessionStepInsertNewComplimentMenuHandler
		kb := tg.NewReplyKeyboardMarkup(
			tg.NewButtonColumn(
				tg.NewKeyboardButton(models.AdminMenu.InsertNewCompliment),
				tg.NewKeyboardButton(models.AdminStartMenu.AdminEnter),
			)...,
		).WithResizeKeyboardMarkup()
		return msg.Update.Reply(ctx, msg.Answer(fmt.Sprintf("Пришлите комплимент, который вы хотите добавить")).
			ReplyMarkup(kb))
	case models.AdminStartMenu.AdminEnter:
		h.sessionManager.Get(ctx).Step = models.SessionStepEnterAdminMenuHandler
		return h.AdminStartMenuHandler(ctx, msg)
	default:
		h.sessionManager.Get(ctx).Step = models.SessionStepInit
		return msg.Answer("Напишите /start").ReplyMarkup(tg.NewReplyKeyboardRemove()).DoVoid(ctx)
	}
}

func (h *AdminHandler) AdminMenuGetAllPhotosHandler(ctx context.Context, msg *tgb.MessageUpdate) error {
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
