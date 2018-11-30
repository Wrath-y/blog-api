package articleController

import (
	"github.com/gin-gonic/gin"
	"go-blog/model/article"
	"go-blog/server/errno"
	"go-blog/struct"
	"go-blog/struct/article-struct"
	"reflect"
	"strconv"
)

func Store(c *gin.Context) {
	var r article_struct.Request
	var err error
	if err := c.Bind(&r); err != nil {
		_struct.Response(c, errno.BindError, nil)

		return
	}

	t := reflect.TypeOf(r)
	v := reflect.ValueOf(r)
	for k := 0; k < t.NumField(); k++ {
		switch t.Field(k).Type.String() {
		case "string":
			if v.Field(k).String() == "" {
				err = errno.New(errno.RequestError, " "+t.Field(k).Name + " can not be null")
				_struct.Response(c, err, nil)

				return
			}
		}
	}

	res := article.Articles{
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

	_struct.Response(c, err, res)

	return
}

func Delete(c *gin.Context) {

}

func Update(c *gin.Context) {

}

func Index(c *gin.Context) {
	page, err:= strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		panic(err)
	}
	data, count, err := article.Index(page, 6)

	if err != nil {
		_struct.Response(c, err, nil)
		return
	}

	_struct.Response(c, nil, article_struct.Response{
		Count: count,
		Data:   data,
	})

	return
}

func Show(c *gin.Context) {

}