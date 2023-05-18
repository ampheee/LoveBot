package user

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/mr-linch/go-tg/tgb/session"
	"tacy/internal/models"
	"tacy/internal/services/userService"
)

type UserHandler struct {
	sessionManager *session.Manager[models.Session]
	service        userService.Service
}

func NewUserHandler(
	sm *session.Manager[models.Session],
	UserService userService.Service,
) *UserHandler {
	return &UserHandler{
		sessionManager: sm,
		service:        UserService,
	}
}

func (h *UserHandler) UserStartMenuSelectionHandler(ctx context.Context, msg *tgb.MessageUpdate) error {
	switch msg.Text {
	case models.UserStartMenu.GetComplimentNow:
		h.sessionManager.Get(ctx).Step = models.SessionStepGetCompliment
		kb := tg.NewReplyKeyboardMarkup(
			tg.NewButtonColumn(
				tg.NewKeyboardButton(models.UserStartMenu.GetComplimentNow),
				tg.NewKeyboardButton(models.UserStartMenu.InsertSomeThoughts),
			)...,
		).WithResizeKeyboardMarkup()
		return msg.Update.Reply(ctx, msg.Answer(fmt.Sprintf("Отлично! Что дальше ?)")).
			ReplyMarkup(kb))
	case models.UserStartMenu.InsertSomeThoughts:
		h.sessionManager.Get(ctx).Step = models.SessionStepInsertSomeThoughts
		kb := tg.NewReplyKeyboardMarkup(
			tg.NewButtonColumn(
				tg.NewKeyboardButton(models.UserStartMenu.GetComplimentNow),
				tg.NewKeyboardButton(models.UserStartMenu.InsertSomeThoughts),
			)...,
		).WithResizeKeyboardMarkup()
		return msg.Update.Reply(ctx, msg.Answer(fmt.Sprintf("Спасибо, уверен, мне будет очень приятно читать:)")).
			ReplyMarkup(kb))
	default:
		h.sessionManager.Get(ctx).Step = models.SessionStepInit
		return msg.Answer("Что-то сломалось:(\nНапиши /start").ReplyMarkup(tg.NewReplyKeyboardRemove()).DoVoid(ctx)
	}
}
