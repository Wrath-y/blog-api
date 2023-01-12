package api

import (
	"blog-api/core"
	"blog-api/errcode"
	"blog-api/server/pixivoss"
	"strconv"
)

func GetPixivs(c *core.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page_size", "15"))
	list, err := pixivoss.List(c.DefaultQuery("next_marker", ""), page)
	if err != nil {
		c.ErrorL("获取pixiv失败", page, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}

	c.Success(list)
}
