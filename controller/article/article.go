package article

import (
	"github.com/gin-gonic/gin"
	"go-blog/model/article"
	"go-blog/req_struct"
	"go-blog/req_struct/req_article"
	"go-blog/server/errno"
	"strconv"
)

func Store(c *gin.Context) {
	var r req_article.Request
	if err := c.Bind(&r); err != nil {
		req_struct.Response(c, errno.BindError, nil)
		return
	}

	if err := r.Validate(c); err != nil {
		req_struct.Response(c, err, nil)
		return
	}

	res := &article.Articles{
		Title: r.Title,
		Image: r.Image,
		Html:  r.Html,
		Con:   r.Con,
		Tags:  r.Tags,
	}
	if err := res.Create(); err != nil {
		req_struct.Response(c, errno.DatabaseError, nil)
		return
	}

	req_struct.Response(c, nil, res)

	return
}

func Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := article.Delete(id); err != nil {
		req_struct.Response(c, errno.DatabaseError, nil)
		return
	}

	req_struct.Response(c, nil, nil)

	return
}

func Update(c *gin.Context) {
	var r req_article.Request
	id, _ := strconv.Atoi(c.Param("id"))
	if err := c.Bind(&r); err != nil {
		req_struct.Response(c, errno.BindError, nil)
		return
	}

	if err := r.Validate(c); err != nil {
		req_struct.Response(c, err, nil)
		return
	}

	res := &article.Articles{
		Title: r.Title,
		Image: r.Image,
		Html:  r.Html,
		Con:   r.Con,
		Tags:  r.Tags,
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
		req_struct.Response(c, err, nil)
		return
	}

	data, count, err := article.Index(page, 15)

	if err != nil {
		req_struct.Response(c, err, nil)
		return
	}

	req_struct.Response(c, nil, req_article.Response{
		Count: count,
		Data:  data,
	})

	return
}

func Show(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := article.Show(id)
	if err != nil {
		req_struct.Response(c, errno.DatabaseError, nil)
		return
	}
	req_struct.Response(c, nil, res)

	return
}
