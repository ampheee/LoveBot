package client

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"tacy/config"
	"tacy/internal/handlers"
	repository2 "tacy/internal/services/authService/repository"
	"tacy/internal/storage/repository"
	"tacy/pkg/botlogger"
)

func Run(ctx context.Context, client *tg.Client, config config.Config, pool *pgxpool.Pool) error {
	logger := botlogger.GetLogger()
	storage := repository.InitStorage(pool)
	auth := repository2.InitNewAuthService(ctx, storage)
	handlers := handlers.New(handlers.Deps{
		AuthService: auth,
		UserService: nil,
	})
	router := handlers.Init(ctx)
	_, err := client.GetMe().Do(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to start bot. Check token and start again")
	}
	return tgb.NewPoller(router, client).Run(ctx)
}
