package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"tacy/config"
	"tacy/pkg/botlogger"
	"tacy/pkg/middleware"
	"time"
)

func GetPool(ctx context.Context, c config.Config) (connect *pgxpool.Pool) {
	logger := botlogger.GetLogger()
	err := middleware.ConnectWithTries(func() error {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		var err error
		connect, err = pgxpool.Connect(ctx, fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?%s&%s",
			c.Postgresql.User,
			c.Postgresql.Pass,
			c.Postgresql.Host,
			c.Postgresql.Port,
			c.Postgresql.DbName,
			c.Postgresql.SSLMode,
			c.Postgresql.MaxConn))
		return err
	}, 3, time.Second*3)
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to connect to database")
	}
	logger.Info().Msg("connected to database successfully")
	return connect
}
