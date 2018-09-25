package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/go-blog/go-blog/handler/health-check"
)

func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine  {
	// middleware
	g.Use(gin.Recovery())
	g.Use(gin.Logger())
	g.Use(mw...)
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "路由错误")
	})

	hCheck := g.Group("/health-check")
	{
		hCheck.GET("/health", health_check.HealthCheck)
		hCheck.GET("/disk", health_check.DiskCheck)
		hCheck.GET("/cpu", health_check.CPUCheck)
		hCheck.GET("/ram", health_check.RAMCheck)
	}

	return g
}