package spiderController

import (
	"github.com/gin-gonic/gin"
	"go-blog/server/spider"
)

func Store(c *gin.Context) {
	spider.Login(c)
}