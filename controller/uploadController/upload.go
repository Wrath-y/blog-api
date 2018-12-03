package uploadController

import (
	"github.com/gin-gonic/gin"
	"go-blog/server/upload"
	"go-blog/struct"
)

func Index(c *gin.Context) {
	data := upload.Version()
	_struct.Response(c, nil, data)

	return
}

func Create(c *gin.Context) {
	upload.Create(c)
}