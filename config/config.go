package config

import (
	"github.com/spf13/viper"
	"tacy/pkg/botlogger"
)

type Config struct {
	Token      string
	Postgresql struct {
		User    string
		Pass    string
		Host    string
		Port    string
		DbName  string
		SSLMode string
		MaxConn string
	}
	AcceptedUser string
}

func LoadConfig() *viper.Viper {
	log := botlogger.GetLogger()
	v := viper.New()
	v.AddConfigPath("../config")
	v.SetConfigName("config")
	v.SetConfigType("yml")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to read config")
	}
	log.Info().Msg("Config loaded successfully")
	return v
}

func ParseConfig(v *viper.Viper) Config {
	var (
		c      Config
		logger = botlogger.GetLogger()
	)
	err := v.Unmarshal(&c)
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to parse config")
	}
	logger.Info().Msg("config parsed successfully")
	return c
}
