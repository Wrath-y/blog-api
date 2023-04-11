package articleseo

import (
	"blog-api/controller/resp"
	"blog-api/core"
	"blog-api/entity"
	"github.com/go-redis/redis/v7"
)

func List(c *core.Context, id int) ([]*resp.GetArticleSeoResp, error) {
	list := make([]*resp.GetArticleSeoResp, 0)
	logMap := make(map[string]interface{})
	logMap["id"] = id
	list, err := GetListById(id)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if list != nil {
		return list, nil
	}

	defer func() {
		if len(list) > 0 {
			_ = SetList(id, list)
		}
		c.Info("从DB获取文章seo列表", logMap, nil)
	}()

	articleSeoList, err := new(entity.ArticleSeo).FindByArticleId(id)
	if err != nil {
		return nil, err
	}

	for _, v := range articleSeoList {
		data := &resp.GetArticleSeoResp{
			Name:    v.Name,
			Content: v.Content,
		}
		list = append(list, data)
	}

	return list, nil
}
