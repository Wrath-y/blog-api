package articleController

import (
	"github.com/gin-gonic/gin"
	"go-blog/model/article"
	"go-blog/server/errno"
	"go-blog/struct"
	"go-blog/struct/articleStruct"
	"strconv"
)

func Store(c *gin.Context) {
	var r articleStruct.Request
	if err := c.Bind(&r); err != nil {
		_struct.Response(c, errno.BindError, nil)
		return
	}

	if err := r.Validate(c); err != nil {
		_struct.Response(c, err, nil)
		return
	}

	res := &article.Articles{
		Title: r.Title,
		Image: r.Image,
		Html: r.Html,
		Con: r.Con,
		Tag: r.Tag,
	}
	if err := res.Create(); err != nil {
		_struct.Response(c, errno.DatabaseError, nil)
		return
	}

	_struct.Response(c, nil, res)

	return
}

func Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := article.Delete(id); err != nil {
		_struct.Response(c, errno.DatabaseError, nil)
		return
	}

	_struct.Response(c, nil, nil)

	return
}

func Update(c *gin.Context) {
	var r articleStruct.Request
	id, _ := strconv.Atoi(c.Param("id"))
	if err := c.Bind(&r); err != nil {
		_struct.Response(c, errno.BindError, nil)
		return
	}

	if err := r.Validate(c); err != nil {
		_struct.Response(c, err, nil)
		return
	}

	res := &article.Articles{
		Title: r.Title,
		Image: r.Image,
		Html: r.Html,
		Con: r.Con,
		Tag: r.Tag,
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
	data, count, err := article.Index(page, 15)

	if err != nil {
		_struct.Response(c, err, nil)
		return
	}

	_struct.Response(c, nil, articleStruct.Response{
		Count: count,
		Data:   data,
	})

	return
}

func Show(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := article.Show(id)
	if err != nil {
		_struct.Response(c, errno.DatabaseError, nil)
		return
	}
	_struct.Response(c, nil, res)

	return
}