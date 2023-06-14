package handlers

import (
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/mr-linch/go-tg/tgb/session"
	"github.com/rs/zerolog"
	"golang.org/x/net/context"
	"tacy/internal/handlers/user"
	"tacy/internal/services/authService"
	"tacy/internal/services/userService"
	"tacy/pkg/botlogger"

	"tacy/internal/handlers/admin"
	"tacy/internal/handlers/start"
	"tacy/internal/models"
	"time"
)

type Deps struct {
	sessionManager *session.Manager[models.Session]
	AuthService    authService.Service
	UserService    userService.USUsecase
	Client         *tg.Client
}

type Handler struct {
	Logger zerolog.Logger
	*tgb.Router
	sessionManager *session.Manager[models.Session]
	StartHandler   *start.StartHandler
	AdminHandler   *admin.AdminHandler
	UserHandler    *user.UserHandler
	Client         *tg.Client
}

func New(deps Deps) *Handler {
	sm := NewSessionManager()
	return &Handler{
		Logger:         botlogger.GetLogger(),
		Router:         tgb.NewRouter(),
		sessionManager: sm.Manager,
		StartHandler:   start.NewStartHandler(sm.Manager, deps.AuthService),
		AdminHandler: admin.NewAdminHandler(
			sm.Manager,
			deps.UserService,
			botlogger.GetLogger(),
			deps.Client,
		),
		UserHandler: user.NewUserHandler(
			sm.Manager,
			deps.UserService,
			botlogger.GetLogger(),
			deps.Client,
		),
	}
}

func (h *Handler) Init(ctx context.Context) *tgb.Router {
	logger := botlogger.GetLogger()
	h.Router.Use(tgb.MiddlewareFunc(func(next tgb.Handler) tgb.Handler {
		return tgb.HandlerFunc(func(ctx context.Context, update *tgb.Update) error {
			defer func(started time.Time) {
				logger.Printf("%#v [%s]", update, time.Since(started))
			}(time.Now())

			return next.Handle(ctx, update)
		})
	}))
	h.registerSession()
	h.registerStartHandlers()
	h.registerAdminHandlers()
	h.registerUserHandler()

	return h.Router
}

func (h *Handler) registerStartHandlers() {
	h.Message(h.StartHandler.Start, tgb.Command("start")).
		Message(func(ctx context.Context, mu *tgb.MessageUpdate) error {
			h.Logger.Warn().Msg(mu.Message.Text + " fetched. Unknown endpoint")
			return mu.Update.Reply(ctx, mu.Answer("Напишите /start"))
		}, h.isSessionStep(models.SessionStepInit))
	h.Logger.Info().Msg("bot started")
}

func (h *Handler) registerAdminHandlers() {
	h.Message(h.AdminHandler.AdminStartMenuHandler, h.isSessionStep(models.SessionStepAdminMenuHandler)).
		Message(h.AdminHandler.AdminMenuHandler, h.isSessionStep(models.SessionStepEnterAdminMenuHandler)).
		Message(h.AdminHandler.AdminMenuInsertPhotoHandler, h.isSessionStep(models.SessionStepInsertNewPhotoMenuHandler)).
		Message(h.AdminHandler.AdminMenuInsertNewComplimentHandler, h.isSessionStep(models.SessionStepInsertNewComplimentMenuHandler)).
		Message(h.AdminHandler.AdminMenuGetAllPhotosHandler, h.isSessionStep(models.SessionStepGetAllPhotosMenuHandler)).
		Message(h.AdminHandler.AdminMenuGetAllComplimentsHandler, h.isSessionStep(models.SessionStepGetAllComplimentsHandler)).
		Message(h.AdminHandler.AdminGetComplimentByRandomHandler, h.isSessionStep(models.SessionStepGetCompliment))
}

func (h *Handler) registerUserHandler() {
	h.Message(h.UserHandler.UserStartMenuSelectionHandler, h.isSessionStep(models.SessionStepUserMenuHandler)).
		Message(h.UserHandler.UserInputSomeThoughts, h.isSessionStep(models.SessionStepInsertSomeThoughts)).
		Message(h.UserHandler.UserGetComplimentByRandomHandler, h.isSessionStep(models.SessionStepGetCompliment))
}
