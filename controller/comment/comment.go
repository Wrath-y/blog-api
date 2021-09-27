package comment

import (
	"github.com/gin-gonic/gin"
	"go-blog/model/comment"
	"go-blog/req_struct"
	"go-blog/req_struct/req_comment"
	"go-blog/server/errno"
	"strconv"
)

func Index(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		panic(err)
	}
	data, count, err := comment.Index(page, 15)
	if err != nil {
		req_struct.Response(c, err, nil)
		return
	}
	req_struct.Response(c, nil, req_comment.Response{
		Count: count,
		Data:  data,
	})
	return
}

func Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := comment.Delete(id); err != nil {
		req_struct.Response(c, errno.DatabaseError, nil)
		return
	}
	req_struct.Response(c, nil, nil)
	return
}
