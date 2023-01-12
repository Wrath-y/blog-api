package core

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const Ctx = "__context__"

type HandlerFunc func(c *Context)

func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := getContext(c)
		h(ctx)
	}
}

func getContext(c *gin.Context) *Context {
	ctx, ok := c.Get(Ctx)
	if ok {
		return ctx.(*Context)
	}

	context := &Context{
		Env:     viper.GetString("app.env"),
		Context: c,
	}
	c.Set(Ctx, context)
	return context
}
