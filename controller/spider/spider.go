package spider

import (
	"github.com/gin-gonic/gin"
	"go-blog/req_struct"
	"go-blog/req_struct/req_spider"
	"go-blog/server/errno"
	"go-blog/server/spider"
	"strconv"
)

func Index(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page_size", "15"))
	list, err := spider.Index(c.DefaultQuery("next_marker", ""), page)
	if err != nil {
		req_struct.Response(c, errno.ServerError, err)
		return
	}
	req_struct.Response(c, nil, list)

	return
}

func Store(c *gin.Context) {
	var r req_spider.UpdateImgRequest
	if err := c.Bind(&r); err != nil {
		req_struct.Response(c, errno.BindError, nil)
		return
	}

	if err := r.Validate(c); err != nil {
		req_struct.Response(c, err, nil)
		return
	}

	spider.Get(c, r.Cookie)
	req_struct.Response(c, nil, nil)

	return
}

func Delete(c *gin.Context) {
	res, err := spider.Delete(c.Query("name"))
	req_struct.Response(c, err, res)
	return
}

func Count(c *gin.Context) {
	res, err := spider.Count()

	req_struct.Response(c, err, res)

	return
}
