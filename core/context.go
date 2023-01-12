package core

import (
	"blog-api/errcode"
	"blog-api/pkg/def"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Logger interface {
	Info(message string, request, response interface{}, t ...time.Time)
	Warn(message string, request, response interface{}, t ...time.Time)
	ErrorL(message string, request, response interface{}, t ...time.Time)
	Fatal(message string, request, response interface{}, t ...time.Time)
}

type Context struct {
	Env string
	*gin.Context
	Logger
	TimeLocation *time.Location
}

func (c *Context) Success(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}

func (c *Context) Fail(code int, msg string, detail, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	if c.Env == def.EnvProduction {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
	} else {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":   code,
			"msg":    msg,
			"detail": detail,
			"data":   data,
		})
	}
}

func (c *Context) FailWithErrCode(err *errcode.ErrCode, data interface{}) {
	c.Fail(err.Code, err.Msg, err.Detail, data)
}
