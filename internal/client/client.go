package client

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"tacy/internal/handlers"
	repository2 "tacy/internal/services/authService/repository"
	repository3 "tacy/internal/services/userService/repository"
	"tacy/internal/storage/repository"
	"tacy/pkg/botlogger"
)

func Run(ctx context.Context, client *tg.Client, pool *pgxpool.Pool) error {
	logger := botlogger.GetLogger()
	storage := repository.InitStorage(pool)
	auth := repository2.InitNewAuthService(storage)
	user := repository3.InitNewUserService(storage)
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
