package router

import (
	"github.com/gin-gonic/gin"
	"go-blog/controller"
	"go-blog/router/middleware"
	"go-blog/server/errno"
)

func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// middleware
	g.Use(gin.Recovery())
	g.Use(middleware.Logger)
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(mw...)
	g.NoRoute(func(c *gin.Context) {
		controller.Response(c, errno.RouteError, nil)
	})

	loadAdmin(g)
	loadApi(g)

	return g
}
