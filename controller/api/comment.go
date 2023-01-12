package api

import (
	"blog-api/core"
	"blog-api/entity"
	"blog-api/errcode"
	"encoding/json"
	"strconv"
	"time"
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

func GetComments(c *core.Context) {
	logMap := make(map[string]interface{})
	lastId, err := strconv.Atoi(c.DefaultQuery("last_id", "1"))
	if err != nil {
		panic(err)
	}
	logMap["lastId"] = lastId
	articleId, err := strconv.Atoi(c.DefaultQuery("article_id", "0"))
	if err != nil {
		panic(err)
	}
	logMap["articleId"] = articleId
	list, err := new(entity.Comment).FindByArticleIdLastId(articleId, lastId, 15)
	if err != nil {
		c.ErrorL("获取评论失败", logMap, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}

	c.Success(list)
}

func AddComment(c *core.Context) {
	var r CommentRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.FailWithErrCode(errcode.WebInvalidParam, nil)
		return
	}

	jsonByte, _ := json.Marshal(r)
	comment := new(entity.Comment)
	err := json.Unmarshal(jsonByte, comment)
	if err != nil {
		c.ErrorL("反序列化失败", r, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}

	comment.CreatedAt = time.Now().In(c.TimeLocation)
	comment.UpdatedAt = time.Now().In(c.TimeLocation)
	if err := comment.Create(); err != nil {
		c.ErrorL("添加评论失败", comment, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}

	c.Success(comment)
}

func GetCommentCount(c *core.Context) {
	articleId, err := strconv.Atoi(c.DefaultQuery("article_id", "0"))
	if err != nil {
		panic(err)
	}
	count, err := new(entity.Comment).GetArticlesWebCommentCount(articleId)
	if err != nil {
		c.ErrorL("获取评论数量失败", articleId, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}

	c.Success(count)
}
