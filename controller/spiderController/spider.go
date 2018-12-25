package spiderController

import (
	"github.com/gin-gonic/gin"
	"go-blog/server/spider"
	"go-blog/struct"
)

func Index(c *gin.Context) {
	spider.Get(c)
	_struct.Response(c, nil, nil)

	return
}