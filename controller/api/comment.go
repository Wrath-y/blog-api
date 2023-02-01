package api

import (
	"blog-api/core"
	"blog-api/entity"
	"blog-api/errcode"
	"blog-api/service/comment"
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
	var (
		logMap = make(map[string]interface{})
	)
	lastId, err := strconv.Atoi(c.DefaultQuery("last_id", "0"))
	if err != nil {
		c.ErrorL("获取lastId失败", nil, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}
	logMap["lastId"] = lastId
	articleId, err := strconv.Atoi(c.DefaultQuery("article_id", "0"))
	if err != nil {
		c.ErrorL("获取articleId失败", nil, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}
	logMap["articleId"] = articleId
	list, err := comment.FindByArticleIdLastId(c, articleId, lastId)
	if err != nil {
		c.ErrorL("获取评论列表失败", nil, err.Error())
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
	data := new(entity.Comment)
	err := json.Unmarshal(jsonByte, data)
	if err != nil {
		c.ErrorL("反序列化失败", r, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}

	data.CreatedAt = time.Now().In(c.TimeLocation)
	data.UpdatedAt = time.Now().In(c.TimeLocation)
	if err := data.Create(); err != nil {
		c.ErrorL("添加评论失败", data, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}

	c.Success(nil)
}

func GetCommentCount(c *core.Context) {
	articleId, err := strconv.Atoi(c.DefaultQuery("article_id", "0"))
	if err != nil {
		panic(err)
	}
	count, err := comment.GetCommentCount(articleId)
	if err != nil {
		c.ErrorL("获取评论数量失败", articleId, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}

	c.Success(map[string]int64{
		"count": count,
	})
}
