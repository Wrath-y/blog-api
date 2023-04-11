package api

import (
	"blog-api/core"
	"blog-api/errcode"
	"blog-api/service/articleseo"
	"strconv"
)

func GetArticleSeo(c *core.Context) {
	var (
		logMap = make(map[string]interface{})
	)

	id, _ := strconv.Atoi(c.Param("id"))
	logMap["id"] = id

	resp, err := articleseo.List(c, id)
	if err != nil {
		c.ErrorL("获取文章seo列表失败", logMap, err.Error())
		c.FailWithErrCode(errcode.WebNetworkBusy, nil)
		return
	}

	c.Success(resp)
}
