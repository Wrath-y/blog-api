package router

import (
	"github.com/gin-gonic/gin"
	"go-blog/controller/admin"
	"go-blog/controller/health_check"
	"go-blog/controller/user"
	"go-blog/router/middleware"
)

func loadAdmin(g *gin.Engine) {
	g.POST("/login", user.Login)

	a := g.Group("/admin")
	a.Use(middleware.Auth())
	{
		hc := a.Group("/health-checks")
		{
			hc.GET("", health_check.HealthCheck)
			hc.GET("/disk", health_check.DiskCheck)
			hc.GET("/cpu", health_check.CPUCheck)
			hc.GET("/ram", health_check.RAMCheck)
		}
		articles := a.Group("/articles")
		{
			articles.POST("", admin.AddArticle)
			articles.DELETE("/:id", admin.DelArticle)
			articles.PUT("/:id", admin.UpdateArticle)
			articles.GET("", admin.GetArticles)
			articles.GET("/:id", admin.GetArticle)
		}
		uploads := a.Group("/uploads")
		{
			uploads.GET("", admin.GetUpload)
		}
		pixivs := a.Group("/pixivs")
		{
			pixivs.GET("", admin.GetPixivs)
			pixivs.GET("count", admin.GetPixivCount)
			pixivs.POST("", admin.AddPixiv)
			pixivs.DELETE("/:id", admin.DelPixiv)
		}
		comments := a.Group("/comments")
		{
			comments.GET("", admin.GetComments)
			comments.DELETE("/:id", admin.DelComment)
		}
		friend := a.Group("/friends")
		{
			friend.POST("", admin.AddFriend)
			friend.DELETE("/:id", admin.DelFriend)
			friend.PUT("/:id", admin.UpdateFriend)
			friend.GET("", admin.GetFriends)
			friend.GET("/:id", admin.GetFriend)
		}
	}
}
