package spiderController

import (
	"github.com/gin-gonic/gin"
	"go-blog/server/errno"
	"go-blog/server/spider"
	"go-blog/struct"
)

func Index(c *gin.Context) {
	spider.Get(c)
	_struct.Response(c, nil, nil)

	return
}

func Store(c *gin.Context) {
	list, err := spider.Index(15)
	if err != nil {
		_struct.Response(c, errno.ServerError, err)
		return
	}
	_struct.Response(c, nil, list)

	return
}