package spiderController

import (
	"github.com/gin-gonic/gin"
	"go-blog/server/spider"
	"go-blog/struct"
)

func Store(c *gin.Context) {
	spider.Login(c)
	_struct.Response(c, nil, nil)

	return
}