package api

import (
	"github.com/gin-gonic/gin"
	"go-blog/controller"
	"go-blog/server/errno"
	"go-blog/server/spider"
	"strconv"
)

type UpdateImgRequest struct {
	Cookie string `json:"cookie" binding:"required"`
}

func GetPixivs(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page_size", "15"))
	list, err := spider.Index(c.DefaultQuery("next_marker", ""), page)
	if err != nil {
		controller.Response(c, errno.ServerError, err)
		return
	}
	controller.Response(c, nil, list)

	return
}
