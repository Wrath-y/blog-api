package _struct

import (
	"github.com/gin-gonic/gin"
	"go-blog/server/errno"
	"net/http"
)

type Rule struct {
	Code	int			`json: "code"`
	Message string		`json: "message"`
	Data	interface{} `json: "data"`
}

func Response(c *gin.Context, err error, data interface{}) {
	code, message := errno.ReturnErr(err)

	c.JSON(http.StatusOK, Rule{
		Code:	code,
		Message:message,
		Data:	data,
	})

	return
}