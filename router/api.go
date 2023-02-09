package router

import (
	"blog-api/controller/api"
	"blog-api/core"
	"blog-api/middleware"
	"github.com/gin-gonic/gin"
)

func loadApi(r *gin.RouterGroup) {
	a := r.Group("/api", core.Handle(middleware.Logging), core.Handle(middleware.TimeLocation))
	{
		a.GET("/pixivs", core.Handle(api.GetPixivs))
		a.GET("/friends", core.Handle(api.GetFriends))

		ars := a.Group("articles")
		{
			ars.GET("", core.Handle(api.GetArticles))
			ars.GET("/a", core.Handle(api.GetAllArticles))
		}
		ar := a.Group("article")
		{
			ar.GET("/:id", core.Handle(api.GetArticle))
		}

		cm := a.Group("comments")
		{
			cm.GET("", core.Handle(api.GetComments))
			cm.GET("/count", core.Handle(api.GetCommentCount))
			cm.POST("", core.Handle(api.AddComment))
		}
	}
}
