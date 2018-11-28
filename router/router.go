package router

import (
	"github.com/gin-gonic/gin"
	"go-blog/controller/articleController"
	"go-blog/controller/healthCheckController"
	"go-blog/router/middleware"
	"go-blog/server/errno"
	"go-blog/struct"
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
		_struct.Response(c, errno.RouteError, nil)
	})

	hc := g.Group("/health-check")
	{
		hc.GET("", healthCheckController.HealthCheck)
		hc.GET("/disk", healthCheckController.DiskCheck)
		hc.GET("/cpu", healthCheckController.CPUCheck)
		hc.GET("/ram", healthCheckController.RAMCheck)
	}
	a := g.Group("article-dao")
	{
		a.POST("", articleController.Store)
		a.DELETE("/:id", articleController.Delete)
		a.PUT("/:id", articleController.Update)
		a.GET("", articleController.Index)
		a.GET("/:id", articleController.Show)
	}

	return g
}

