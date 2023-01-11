package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-blog/controller"
	"go-blog/entity/comment"
	"go-blog/server/errno"
	"strconv"
)

type CommentRequest struct {
	ArticleId int    `json:"article_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Pid       int    `json:"pid"`
	Ppid      int    `json:"ppid"`
	Url       string `json:"url" binding:"required"`
}

func GetComments(c *gin.Context) {
	lastId, err := strconv.Atoi(c.DefaultQuery("last_id", "1"))
	if err != nil {
		panic(err)
	}
	articleId, err := strconv.Atoi(c.DefaultQuery("article_id", "0"))
	if err != nil {
		panic(err)
	}
	data, err := comment.IndexBuyArticleId(articleId, lastId, 15)
	if err != nil {
		controller.Response(c, err, nil)
		return
	}
	controller.Response(c, nil, data)
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

func AddComment(c *gin.Context) {
	var r CommentRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		controller.Response(c, errno.BindError, err)
		return
	}

	jsonByte, err := json.Marshal(r)
	if err != nil {
		controller.Response(c, errno.BindError, err)
		return
	}
	res := new(comment.Comment)
	err = json.Unmarshal(jsonByte, res)
	if err != nil {
		controller.Response(c, errno.BindError, err)
		return
	}

	if err := res.Create(); err != nil {
		controller.Response(c, errno.DatabaseError, err)
		return
	}

	controller.Response(c, nil, res)
	return
}

func GetCommentCount(c *gin.Context) {
	articleId, err := strconv.Atoi(c.DefaultQuery("article_id", "0"))
	if err != nil {
		panic(err)
	}
	data, err := comment.GetArticlesWebCommentCount(articleId)
	if err != nil {
		controller.Response(c, err, nil)
		return
	}
	controller.Response(c, nil, data.CommentCount)
	return
}
