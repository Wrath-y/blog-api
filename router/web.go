package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine  {
	// middleware
	g.Use(gin.Recovery())
	g.Use(gin.Logger())
	g.Use(mw...)
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "路由错误")
	})

	hCheck := g.Group("/check")
	{
		hCheck.GET("/health", )
	}
}