package api

import (
	"blog-api/core"
	"blog-api/errcode"
	"blog-api/service/friendlink"
)

func GetFriends(c *core.Context) {
	list, err := friendlink.FindAll()
	if err != nil {
		c.ErrorL("获取友链失败", nil, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}

	c.Success(list)
}
