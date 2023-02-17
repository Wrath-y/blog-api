package api

import (
	"blog-api/core"
	"blog-api/errcode"
	"blog-api/service/article"
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

	resp, err := article.List(c, lastId)
	if err != nil {
		c.ErrorL("获取文章列表失败", logMap, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}

	c.Success(resp)
}

func GetAllArticles(c *core.Context) {
	resp, err := article.All(c)
	if err != nil {
		c.ErrorL("获取所有文章失败", nil, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}

	c.Success(resp)
}

func GetArticle(c *core.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := article.Get(id)
	if err != nil {
		c.ErrorL("获取文章失败", id, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}

	c.Success(res)
}
