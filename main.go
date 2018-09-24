package go_blog

import "github.com/gin-gonic/gin"

func main() {
	g := gin.New()
	g.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		});
	})
	g.Run()
}