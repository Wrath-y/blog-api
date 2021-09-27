package harem

import (
	"github.com/gin-gonic/gin"
	"go-blog/model/harem"
	"go-blog/req_struct"
	"go-blog/req_struct/req_harem"
	"go-blog/server/errno"
	"strconv"
)

func Store(c *gin.Context) {
	var r req_harem.Request
	if err := c.Bind(&r); err != nil {
		req_struct.Response(c, errno.BindError, err)
		return
	}

	res := &harem.Harem{
		Name:  r.Name,
		Email: r.Email,
		Url:   r.Url,
	}
	if err := res.Create(); err != nil {
		req_struct.Response(c, errno.DatabaseError, err)
		return
	}

	req_struct.Response(c, nil, res)

	return
}

func Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := harem.Delete(id); err != nil {
		req_struct.Response(c, errno.DatabaseError, nil)
		return
	}

	req_struct.Response(c, nil, nil)

	return
}

func Update(c *gin.Context) {
	var r req_harem.Request
	id, _ := strconv.Atoi(c.Param("id"))
	if err := c.Bind(&r); err != nil {
		req_struct.Response(c, errno.BindError, nil)
		return
	}

	res := &harem.Harem{
		Name:  r.Name,
		Email: r.Email,
		Url:   r.Url,
	}
	if err := res.Update(id); err != nil {
		req_struct.Response(c, errno.DatabaseError, nil)
		return
	}

	req_struct.Response(c, nil, res)

	return
}

func Index(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		panic(err)
	}
	data, count, err := harem.Index(page, 15)

	if err != nil {
		req_struct.Response(c, err, nil)
		return
	}

	req_struct.Response(c, nil, req_harem.Response{
		Count: count,
		Data:  data,
	})

	return
}

func Show(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := harem.Show(id)
	if err != nil {
		req_struct.Response(c, errno.DatabaseError, nil)
		return
	}
	req_struct.Response(c, nil, res)

	return
}
