package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ip2location/ip2location-go/v9"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go-blog/config"
	"go-blog/model"
	"go-blog/router"
	"runtime"
)

func main() {
	var err error
	cpuNum := runtime.NumCPU() - 1
	if cpuNum <= 0 {
		cpuNum = 1
	}
	runtime.GOMAXPROCS(cpuNum)
	pflag.Parse()
	if err := config.Init("/var/www/go-blog/config/config.yaml"); err != nil {
		log.Fatal().Err(err).Msg("init config faild")
	}

	model.DB.Init()
	defer model.DB.Close()

	gin.SetMode(viper.GetString("runmode"))
	g := gin.New()

	router.Load(g)

	// IP转区号
	config.IP2LocationDB, err = ip2location.OpenDB("./IP2LOCATION-LITE-DB1.BIN")
	defer config.IP2LocationDB.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("./IP2LOCATION-LITE-DB1.BIN open faild")
	}

	// g.RunTLS(viper.GetString("port"), viper.GetString("fullchain"), viper.GetString("key"))
	g.Run(":" + viper.GetString("port"))
}
