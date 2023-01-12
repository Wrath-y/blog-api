package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go-blog/pkg/def"
	"log"
)

const DefaultRelationPath = "./conf.yaml"

func Setup() {
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	// listen and auto reload config
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()

	env := viper.GetString("app.env")
	if env != def.EnvDevelopment && env != def.EnvTesting && env != def.StagingEnv && env != def.EnvProduction {
		log.Fatal("app.env is not correct")
	}
}
