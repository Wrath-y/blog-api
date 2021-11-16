package api

import (
	"github.com/gin-gonic/gin"
	"go-blog/controller"
	"go-blog/model/harem"
	"strconv"
)

func GetFriends(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		panic(err)
	}
	data, err := harem.Index(page, 15)

	if err != nil {
		controller.Response(c, err, nil)
		return
	}

	controller.Response(c, nil, data)

	return
}
