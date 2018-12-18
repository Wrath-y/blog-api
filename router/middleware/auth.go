package middleware

import (
	"github.com/gin-gonic/gin"
	"go-blog/server/errno"
	"go-blog/server/token"
	"go-blog/struct"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := token.ParseRequest(c); err != nil {
			_struct.Response(c, errno.ErrTokenInvalid, err)
			c.Abort()

			return
		}

		c.Next()
	}
}