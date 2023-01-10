package router

import (
	"github.com/gin-gonic/gin"
	"go-blog/controller"
	middleware2 "go-blog/middleware"
	"go-blog/server/errno"
)

func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// middleware
	g.Use(gin.Recovery())
	g.Use(middleware2.Logger)
	g.Use(middleware2.NoCache)
	g.Use(middleware2.Options)
	g.Use(mw...)
	g.NoRoute(func(c *gin.Context) {
		controller.Response(c, errno.RouteError, nil)
	})

	loadAdmin(g)
	loadApi(g)

	return g
}
