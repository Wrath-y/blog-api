package spiderController

import (
	"github.com/gin-gonic/gin"
	"go-blog/server/errno"
	"go-blog/server/spider"
	"go-blog/struct"
	"go-blog/struct/spiderStruct"
	"strconv"
)

func Index(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page_size", "15"))
	list, err := spider.Index(c.DefaultQuery("next_marker", ""), page)
	if err != nil {
		_struct.Response(c, errno.ServerError, err)
		return
	}
	_struct.Response(c, nil, list)

	return
}

func Store(c *gin.Context) {
	var r spiderStruct.UpdateImgRequest
	if err := c.Bind(&r); err != nil {
		_struct.Response(c, errno.BindError, nil)
		return
	}

	if err := r.Validate(c); err != nil {
		_struct.Response(c, err, nil)
		return
	}

	spider.Get(c, r.Cookie)
	_struct.Response(c, nil, nil)

	return
}

func Delete(c *gin.Context) {
	res, err := spider.Delete(c.Query("name"))
	_struct.Response(c, err, res)
	return
}

func Count(c *gin.Context) {
	res, err := spider.Count()

	_struct.Response(c, err, res)

	return
}