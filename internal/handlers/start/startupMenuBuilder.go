package start

import (
	"github.com/mr-linch/go-tg"
	"tacy/internal/models"
)

func buildAdminStartMenu() *tg.ReplyKeyboardMarkup {
	return tg.NewReplyKeyboardMarkup(
		tg.NewButtonColumn(
			tg.NewKeyboardButton(models.AdminStartMenu.AdminEnter),
		)...,
	).WithResizeKeyboardMarkup()
}

func buildUserStartMenu() *tg.ReplyKeyboardMarkup {
	return tg.NewReplyKeyboardMarkup(
		tg.NewButtonColumn(
			tg.NewKeyboardButton(models.UserStartMenu.GetComplimentNow),
			tg.NewKeyboardButton(models.UserStartMenu.InsertSomeThoughts),
		)...,
	).WithResizeKeyboardMarkup()
}
