package api

import (
	"github.com/gin-gonic/gin"
	"go-blog/controller"
	"go-blog/model/friend"
)

func GetFriends(c *gin.Context) {
	data, err := friend.WebIndex()

	if err != nil {
		controller.Response(c, err, nil)
		return
	}

	controller.Response(c, nil, data)

	return
}
