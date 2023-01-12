package api

import (
	resp2 "blog-api/controller/resp"
	"blog-api/core"
	"blog-api/entity"
	"blog-api/errcode"
	"strconv"
)

func GetArticles(c *core.Context) {
	var (
		logMap = make(map[string]interface{})
	)

	lastIdStr := c.DefaultQuery("last_id", "0")
	logMap["lastIdStr"] = lastIdStr

	lastId, err := strconv.Atoi(lastIdStr)
	if err != nil {
		c.ErrorL("转换格式失败", logMap, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}

	articles, err := new(entity.Article).FindByLastId(lastId, 6)
	if err != nil {
		c.ErrorL("获取文章列表失败", logMap, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}
	logMap["articles"] = articles
	resp := make([]*resp2.GetArticlesResp, 0, len(articles))

	articleIds := make([]int, 0, len(articles))
	for _, v := range articles {
		articleIds = append(articleIds, v.Id)
	}
	logMap["articleIds"] = articleIds

	commentCounts, err := new(entity.Comment).GetArticlesWebCommentCounts(articleIds)
	if err != nil {
		c.ErrorL("获取评论失败", logMap, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}

	articleCommentCountMap := make(map[int]int)
	for _, v := range commentCounts {
		articleCommentCountMap[v.ArticleId] = v.CommentCount
	}

	for _, v := range articles {
		data := &resp2.GetArticlesResp{
			Article: v,
		}
		if articleCommentCount, ok := articleCommentCountMap[v.Id]; ok {
			data.CommentCount = articleCommentCount
		}
	}

	c.Success(resp)
}

func GetArticle(c *core.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := new(entity.Article).GetById(id)
	if err != nil {
		c.ErrorL("获取文章失败", id, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}

	c.Success(res)
}
