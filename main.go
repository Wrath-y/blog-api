package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go-blog/config"
	"go-blog/model"
	"go-blog/router"
	"runtime"
)

var (
	cfg = pflag.StringP("config", "c", "", "go-blog config file path")
)

func main() {
	cpuNum := runtime.NumCPU() - 1
	if cpuNum <= 0 {
		cpuNum = 1
	}
	runtime.GOMAXPROCS(cpuNum)
	pflag.Parse()
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	model.DB.Init()
	defer model.DB.Close()

	gin.SetMode(viper.GetString("runmode"))
	g := gin.New()

	middlewares := []gin.HandlerFunc{}

	router.Load(
			g,
			middlewares...,
		)

	// g.RunTLS(viper.GetString("port"), viper.GetString("fullchain"), viper.GetString("key"))
	g.Run(viper.GetString("port"))
}
