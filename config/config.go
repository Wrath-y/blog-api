package config

import (
	"github.com/ip2location/ip2location-go/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go-blog/server/logger"
	"os"
	"strings"
)

type Config struct {
	Name string
}

var IP2LocationDB *ip2location.DB

func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}

	if err := c.initConfig(); err != nil {
		return err
	}

	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name)
	} else {
		viper.AddConfigPath("config")
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("BLOG")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// 初始化日志
	_, pathErr := os.Stat(viper.GetString("logPath"))
	if pathErr != nil {
		makeErr := os.MkdirAll(viper.GetString("logPath"), 0755)
		if makeErr != nil {
			log.Fatal().Err(makeErr).Msg("server log path make error")
		}
	}
	zerolog.SetGlobalLevel(zerolog.Level(viper.GetInt("logLevel")))
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	log.Logger = zerolog.New(logger.NewFileWriter(viper.GetString("logPath"), "blog")).With().Timestamp().Logger()

	return nil
}
