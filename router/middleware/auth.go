package middleware

import (
	"github.com/gin-gonic/gin"
	"go-blog/req_struct"
	"go-blog/server/errno"
	"go-blog/server/token"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := token.ParseRequest(c); err != nil {
			req_struct.Response(c, errno.TokenInvalidErr, err)
			c.Abort()

			return
		}

		c.Next()
	}
}
