package config

import (
	"errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go-blog/server/logger"
	"os"
	"strings"
)

type Config struct {
	Path string
}

func Init(path string) error {
	c := Config{
		Path: path,
	}

	if err := c.initConfig(); err != nil {
		return err
	}

	return nil
}

func (c *Config) initConfig() error {
	if c.Path == "" {
		return errors.New("未设置配置文件路径")
	}
	viper.AddConfigPath(c.Path)
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
