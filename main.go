package go_blog

import (
	"github.com/gin-gonic/gin"
	"github.com/go-blog/go-blog/router"
)

func main() {
	g := gin.New()

	middlewares := []gin.HandlerFunc{}

	router.Load(
			g,
			middlewares...,
		)

	g.Run()
}