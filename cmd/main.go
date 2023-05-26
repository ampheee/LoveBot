package main

import (
	"context"
	"github.com/mr-linch/go-tg"
	"tacy/config"
	"tacy/internal/client"
	"tacy/pkg/botlogger"
	"tacy/pkg/postgresql"
)

func main() {
	logger := botlogger.GetLogger()
	ctx := context.Background()
	v := config.LoadConfig()
	Config := config.ParseConfig(v)
	bot := tg.New(Config.Token)
	Pool := postgresql.GetPool(ctx, Config)
	logger.Fatal().Err(client.Run(ctx, bot, Pool))
}
