package client

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"strconv"
	"tacy/config"
	"tacy/internal/handlers"
	repository2 "tacy/internal/services/authService/repository"
	"tacy/internal/services/userService"
	repository3 "tacy/internal/services/userService/repository"
	"tacy/internal/services/userService/usecase"
	"tacy/pkg/botlogger"
	"time"
)

func Run(ctx context.Context, client *tg.Client, pool *pgxpool.Pool, c config.Config) error {
	logger := botlogger.GetLogger()
	auth := repository2.NewAuthUsecase(pool)
	user := usecase.NewUSUsecase(repository3.NewUSRepo(pool))
	ticker := time.NewTicker(6 * time.Hour)
	go Sender(ctx, ticker, user, client, c)
	handlers := handlers.New(handlers.Deps{
		AuthService: auth,
		UserService: user,
		Client:      client,
	})
	router := handlers.Init(ctx)
	_, err := client.GetMe().Do(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to start bot. Check token and start again")
	}
	return tgb.NewPoller(router, client).Run(ctx)
}

func Sender(ctx context.Context, ticker *time.Ticker, user userService.USUsecase, client *tg.Client,
	config config.Config) {
	logger := botlogger.GetLogger()
	for range ticker.C {
		photo, compliment, err := user.OutputComplimentAndPhotoByRandom(ctx)
		if err != nil {
			logger.Warn().Err(err)
			return
		}
		id, err := strconv.ParseInt(config.AcceptedUser, 10, 64)
		user := tg.User{
			ID: tg.UserID(id),
		}
		err = client.SendPhoto(user.ID, tg.NewFileArgUpload(tg.NewInputFileBytes("photo", photo))).DoVoid(ctx)
		if err != nil {
			logger.Warn().Err(err).Msg("Wasn`t send")
		}
		err = client.SendMessage(user.ID, compliment).DoVoid(ctx)
		if err != nil {
			logger.Warn().Err(err).Msg("Wasn`t send")
		}
		if err == nil {
			logger.Info().Msg("send")
		}
	}
}
