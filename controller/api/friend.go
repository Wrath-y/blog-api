package api

import (
	"blog-api/core"
	"blog-api/entity"
	"blog-api/errcode"
)

func GetFriends(c *core.Context) {
	data, err := new(entity.FriendLink).FindAll()
	if err != nil {
		c.ErrorL("获取友链失败", nil, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}

	c.Success(data)
}
