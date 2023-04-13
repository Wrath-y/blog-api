package api

import (
	"blog-api/core"
	"blog-api/entity"
	"blog-api/errcode"
	"blog-api/service/article"
	"blog-api/service/comment"
	"strconv"
	"time"
)

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

type AddCommentRequest struct {
	LastId    int    `json:"last_id" binding:"required"`
	ArticleId int    `json:"article_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Pid       int    `json:"pid"`
	Url       string `json:"url"`
}

func AddComment(c *core.Context) {
	var r *AddCommentRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.FailWithErrCode(errcode.WebInvalidParam, nil)
		return
	}

	if r.Name == "" {
		r.Name = "匿名用户"
	}
	data := &entity.Comment{
		Base: &entity.Base{
			UpdatedAt: time.Now().In(c.TimeLocation),
			CreatedAt: time.Now().In(c.TimeLocation),
		},
		Name:      r.Name,
		Email:     r.Email,
		Url:       r.Url,
		Content:   r.Content,
		ArticleId: r.ArticleId,
		Pid:       r.Pid,
	}

	if err := data.Create(); err != nil {
		c.ErrorL("添加评论失败", data, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}

	if err := comment.ClearCommentCache(r.ArticleId, r.LastId); err != nil {
		c.ErrorL("删除评论缓存失败", data, err.Error())
	}

	if err := article.CacheCommentCountIncr(r.ArticleId); err != nil {
		c.ErrorL("评论数量缓存incr失败", data, err.Error())
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
