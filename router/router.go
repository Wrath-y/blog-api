package router

import (
	"github.com/gin-gonic/gin"
	"go-blog/model/article"
	"net/http"
	"go-blog/controller/healthCheckController"
	"go-blog/router/middleware"
)

func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine  {
	// middleware
	g.Use(gin.Recovery())
	g.Use(gin.Logger())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "路由错误")
	})

	hCheck := g.Group("/healthCheckController")
	{
		hCheck.GET("/health", healthCheckController.HealthCheck)
		hCheck.GET("/disk", healthCheckController.DiskCheck)
		hCheck.GET("/cpu", healthCheckController.CPUCheck)
		hCheck.GET("/ram", healthCheckController.RAMCheck)
	}
	a := g.Group("articles")
	{
		a.POST("", article.Create)
	}

	return g
}

