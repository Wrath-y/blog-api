package admin

import (
	"github.com/gin-gonic/gin"
	"go-blog/controller"
	"go-blog/entity/comment"
	"go-blog/server/errno"
	"strconv"
)

func GetComments(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		panic(err)
	}
	data, count, err := comment.AdminIndex(page, 15)
	if err != nil {
		controller.Response(c, errno.DatabaseError, err)
		return
	}
	controller.Response(c, nil, map[string]interface{}{
		"list":  data,
		"count": count,
	})
	return
}

func DelComment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := comment.Delete(id); err != nil {
		controller.Response(c, errno.DatabaseError, nil)
		return
	}
	controller.Response(c, nil, nil)
	return
}
