package router

import (
	"github.com/gin-gonic/gin"
	"go-blog/controller/api"
)

func loadApi(g *gin.Engine) {
	g.GET("/pixivs", api.GetPixivs)

	h := g.Group("friends")
	{
		h.GET("", api.GetFriends)
	}

	a := g.Group("articles")
	{
		a.GET("", api.GetArticles)
		a.GET("/:id", api.GetArticle)
	}

	c := g.Group("comments")
	{
		c.GET("", api.GetComments)
		c.GET("/count", api.GetCommentCount)
		c.POST("", api.AddComment)
	}
}
