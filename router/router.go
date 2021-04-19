package router

import (
	"github.com/gin-gonic/gin"
	"go-blog/controller/CommentController"
	"go-blog/controller/HaremController"
	"go-blog/controller/articleController"
	"go-blog/controller/healthCheckController"
	"go-blog/controller/spiderController"
	"go-blog/controller/uploadController"
	"go-blog/controller/userController"
	"go-blog/router/middleware"
	"go-blog/server/errno"
	"go-blog/struct"
)

func Load(g *gin.Engine) *gin.Engine  {
	// middleware
	g.Use(gin.Recovery())
	g.Use(middleware.Logger)
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)

	g.NoRoute(func(c *gin.Context) {
		_struct.Response(c, errno.RouteError, nil)
	})

	g.GET("/pixivs", spiderController.Index)

	g.POST("/login", userController.Login)

	harems := g.Group("harems")
	{
		harems.GET("", HaremController.Index)
	}

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
		pixivs := admin.Group("pixivs")
		{
			pixivs.GET("", spiderController.Index)
			pixivs.GET("count", spiderController.Count)
			pixivs.POST("", spiderController.Store)
			pixivs.DELETE("/:id", spiderController.Delete)
		}
		comments := admin.Group("comments")
		{
			comments.GET("", CommentController.Index)
			comments.DELETE("/:id", CommentController.Delete)
		}
		harems := admin.Group("harems")
		{
			harems.POST("", HaremController.Store)
			harems.DELETE("/:id", HaremController.Delete)
			harems.PUT("/:id", HaremController.Update)
			harems.GET("", HaremController.Index)
			harems.GET("/:id", HaremController.Show)
		}
	}

	return g
}

