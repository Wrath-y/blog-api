package main

import (
	"go-blog/router"
	"go-blog/config"
	"go-blog/model"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg = pflag.StringP("config", "c", "", "go-blog config file path")
)

func main() {
	pflag.Parse()
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	model.DB.Init()
	defer model.DB.Close()

	g := gin.New()
	gin.SetMode(viper.GetString("runmode"))

	middlewares := []gin.HandlerFunc{}

	router.Load(
			g,
			middlewares...,
		)

	g.Run(viper.GetString("port"))
}