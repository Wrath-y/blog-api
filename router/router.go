package router

import (
	"blog-api/core"
	"blog-api/errcode"
	"blog-api/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register() *gin.Engine {
	r := gin.New()
	r.Use(middleware.Recovery)
	r.Use(middleware.SetV())
	r.Use(core.Handle(middleware.CORS))
	r.NoRoute(NoRoute)

	r.Any("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	g := r.Group("/")
	loadApi(g)

	return r
}

func NoRoute(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotFound, errcode.LibNoRoute)
}
