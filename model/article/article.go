package article

import (
	"github.com/gin-gonic/gin"
	"go-blog/server/errno"
	"go-blog/struct"
	"go-blog/struct/article-struct"
	"reflect"
)

func Create(c *gin.Context) {
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

	res := article_struct.Response{
		Title: r.Title,
		Image: r.Image,
		Html: r.Html,
		Con: r.Con,
	}

	_struct.Response(c, nil, res)

	return
}