package router

import (
	"github.com/gin-gonic/gin"
	"go-blog/controller/articleController"
	"go-blog/controller/healthCheckController"
	"go-blog/controller/uploadController"
	"go-blog/controller/userController"
	"go-blog/router/middleware"
	"go-blog/server/errno"
	"go-blog/struct"
)

func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine  {
	// middleware
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(mw...)
	g.NoRoute(func(c *gin.Context) {
		_struct.Response(c, errno.RouteError, nil)
	})

	g.POST("/login", userController.Login)

	admin := g.Group("/admin")
	admin.Use(middleware.Auth())
	{
		hc := admin.Group("/health-checks")
		{
			hc.GET("", healthCheckController.HealthCheck)
			hc.GET("/disk", healthCheckController.DiskCheck)
			hc.GET("/cpu", healthCheckController.CPUCheck)
			hc.GET("/ram", healthCheckController.RAMCheck)
		}
		articles := admin.Group("articles")
		{
			articles.POST("", articleController.Store)
			articles.DELETE("/:id", articleController.Delete)
			articles.PUT("/:id", articleController.Update)
			articles.GET("", articleController.Index)
			articles.GET("/:id", articleController.Show)
		}
		uploads := admin.Group("uploads")
		{
			uploads.GET("", uploadController.Index)
		}
	}

	return g
}

