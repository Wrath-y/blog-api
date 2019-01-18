package HaremController

import (
	"github.com/gin-gonic/gin"
	"go-blog/model/harem"
	"go-blog/server/errno"
	"go-blog/struct"
	"go-blog/struct/haremStruct"
	"strconv"
)

func Store(c *gin.Context) {
	var r haremStruct.Request
	if err := c.Bind(&r); err != nil {
		_struct.Response(c, errno.BindError, err)
		return
	}

	res := &harem.Harem{
		Name: r.Name,
		Email: r.Email,
		Url: r.Url,
	}
	if err := res.Create; err != nil {
		_struct.Response(c, errno.DatabaseError, err())
		return
	}

	_struct.Response(c, nil, res)

	return
}

func Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := harem.Delete(id); err != nil {
		_struct.Response(c, errno.DatabaseError, nil)
		return
	}

	_struct.Response(c, nil, nil)

	return
}

func Update(c *gin.Context) {
	var r haremStruct.Request
	id, _ := strconv.Atoi(c.Param("id"))
	if err := c.Bind(&r); err != nil {
		_struct.Response(c, errno.BindError, nil)
		return
	}

	res := &harem.Harem{
		Name: r.Name,
		Email: r.Email,
		Url: r.Url,
	}
	if err := res.Update(id); err != nil {
		_struct.Response(c, errno.DatabaseError, nil)
		return
	}

	_struct.Response(c, nil, res)

	return
}

func Index(c *gin.Context) {
	page, err:= strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		panic(err)
	}
	data, count, err := harem.Index(page, 15)

	if err != nil {
		_struct.Response(c, err, nil)
		return
	}

	_struct.Response(c, nil, haremStruct.Response{
		Count: count,
		Data:   data,
	})

	return
}