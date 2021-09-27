package router

import (
	"github.com/gin-gonic/gin"
	"go-blog/controller/article"
	"go-blog/controller/comment"
	"go-blog/controller/harem"
	"go-blog/controller/health_check"
	"go-blog/controller/spider"
	"go-blog/controller/upload"
	"go-blog/controller/user"
	"go-blog/req_struct"
	"go-blog/router/middleware"
	"go-blog/server/errno"
)

func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// middleware
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(mw...)
	g.NoRoute(func(c *gin.Context) {
		req_struct.Response(c, errno.RouteError, nil)
		return
	})

	g.GET("/pixivs", spider.Index)

	g.POST("/login", user.Login)

	harems := g.Group("harems")
	{
		harems.GET("", harem.Index)
	}

	admin := g.Group("/admin")
	admin.Use(middleware.Auth())
	{
		hc := admin.Group("/health-checks")
		{
			hc.GET("", health_check.HealthCheck)
			hc.GET("/disk", health_check.DiskCheck)
			hc.GET("/cpu", health_check.CPUCheck)
			hc.GET("/ram", health_check.RAMCheck)
		}
		articles := admin.Group("articles")
		{
			articles.POST("", article.Store)
			articles.DELETE("/:id", article.Delete)
			articles.PUT("/:id", article.Update)
			articles.GET("", article.Index)
			articles.GET("/:id", article.Show)
		}
		uploads := admin.Group("uploads")
		{
			uploads.GET("", upload.Index)
		}
		pixivs := admin.Group("pixivs")
		{
			pixivs.GET("", spider.Index)
			pixivs.GET("count", spider.Count)
			pixivs.POST("", spider.Store)
			pixivs.DELETE("/:id", spider.Delete)
		}
		comments := admin.Group("comments")
		{
			comments.GET("", comment.Index)
			comments.DELETE("/:id", comment.Delete)
		}
		harems := admin.Group("harems")
		{
			harems.POST("", harem.Store)
			harems.DELETE("/:id", harem.Delete)
			harems.PUT("/:id", harem.Update)
			harems.GET("", harem.Index)
			harems.GET("/:id", harem.Show)
		}
	}

	return g
}
