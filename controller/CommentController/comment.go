package CommentController

import (
	"github.com/gin-gonic/gin"
	"go-blog/model/comment"
	"go-blog/server/errno"
	"go-blog/struct"
	"go-blog/struct/commentStruct"
	"strconv"
)

func Index(c *gin.Context) {
	page, err:= strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		panic(err)
	}
	data, count, err := comment.Index(page, 15)
	if err != nil {
		_struct.Response(c, err, nil)
		return
	}
	_struct.Response(c, nil, commentStruct.Response{
		Count: count,
		Data:   data,
	})
	return
}

func Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := comment.Delete(id); err != nil {
		_struct.Response(c, errno.DatabaseError, nil)
		return
	}
	_struct.Response(c, nil, nil)
	return
}