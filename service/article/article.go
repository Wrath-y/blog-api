package article

import (
	"blog-api/controller/resp"
	"blog-api/core"
	"blog-api/entity"
	"github.com/go-redis/redis/v7"
	"gorm.io/gorm"
)

func List(c *core.Context, lastId int) ([]*resp.GetArticlesResp, error) {
	list := make([]*resp.GetArticlesResp, 0)
	logMap := make(map[string]interface{})
	logMap["lastId"] = lastId
	list, err := GetListByLastId(lastId)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if list != nil {
		return list, nil
	}

	defer func() {
		if len(list) > 0 {
			_ = SetList(lastId, list)
		}
		c.Info("从DB获取文章列表", logMap, nil)
	}()

	articles, err := new(entity.Article).FindByLastId(lastId, 6)
	if err != nil {
		return nil, err
	}
	logMap["articlesLen"] = len(articles)

	articleIds := make([]int, 0, len(articles))
	for _, v := range articles {
		articleIds = append(articleIds, v.Id)
	}
	logMap["articleIds"] = articleIds

	commentCounts, err := new(entity.Comment).GetArticlesWebCommentCounts(articleIds)
	if err != nil {
		c.ErrorL("获取评论失败", logMap, err.Error())
		return nil, err
	}
	logMap["commentCounts"] = commentCounts

	articleCommentCountMap := make(map[int]int64)
	for _, v := range commentCounts {
		articleCommentCountMap[v.ArticleId] = v.CommentCount
	}

	for _, v := range articles {
		data := &resp.GetArticlesResp{
			Article: v,
		}
		if articleCommentCount, ok := articleCommentCountMap[v.Id]; ok {
			data.CommentCount = articleCommentCount
		}
		list = append(list, data)
	}

	return list, nil
}

func All(c *core.Context, lastId int) ([]*resp.GetArticlesResp, error) {
	list := make([]*resp.GetArticlesResp, 0)
	logMap := make(map[string]interface{})
	logMap["lastId"] = lastId
	list, err := GetListByLastId(lastId)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if list != nil {
		return list, nil
	}

	defer func() {
		if len(list) > 0 {
			_ = SetList(lastId, list)
		}
		c.Info("从DB获取文章列表", logMap, nil)
	}()

	articles, err := new(entity.Article).FindByLastId(lastId, 6)
	if err != nil {
		return nil, err
	}
	logMap["articlesLen"] = len(articles)

	articleIds := make([]int, 0, len(articles))
	for _, v := range articles {
		articleIds = append(articleIds, v.Id)
	}
	logMap["articleIds"] = articleIds

	commentCounts, err := new(entity.Comment).GetArticlesWebCommentCounts(articleIds)
	if err != nil {
		c.ErrorL("获取评论失败", logMap, err.Error())
		return nil, err
	}
	logMap["commentCounts"] = commentCounts

	articleCommentCountMap := make(map[int]int64)
	for _, v := range commentCounts {
		articleCommentCountMap[v.ArticleId] = v.CommentCount
	}

	for _, v := range articles {
		data := &resp.GetArticlesResp{
			Article: v,
		}
		if articleCommentCount, ok := articleCommentCountMap[v.Id]; ok {
			data.CommentCount = articleCommentCount
		}
		list = append(list, data)
	}

	return list, nil
}

func Get(id int) (*resp.GetArticlesResp, error) {
	res, err := GetById(id)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if res != nil && res.Id > 0 {
		return res, nil
	}

	article, err := new(entity.Article).GetById(id)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = Set(id, res)
	}()
	comment, err := new(entity.Comment).GetArticlesWebCommentCount(article.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	res = &resp.GetArticlesResp{
		Article:      article,
		CommentCount: comment,
	}

	return res, nil
}
