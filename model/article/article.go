package article

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-blog/server/errno"
	"go-blog/struct"
	"go-blog/struct/article-struct"
)

func Create(c *gin.Context) {
	var r article_struct.Request
	var err error
	if err := c.Bind(&r); err != nil {
		_struct.Response(c, errno.BindError, nil)

		return
	}

	if r.Title == "" {
		err = errno.New(errno.TitleError, fmt.Errorf("title can not be null"))
		_struct.Response(c, err, nil)

		return
	}

	res := article_struct.Response{
		Title: r.Title,
	}

	_struct.Response(c, nil, res)

	return
}