package article

import (
	resp2 "blog-api/controller/resp"
	"blog-api/core"
	"blog-api/entity"
	"github.com/go-redis/redis/v7"
	"gorm.io/gorm"
)

func List(c *core.Context, lastId int) ([]*resp2.GetArticlesResp, error) {
	resp, err := GetListByLastId(lastId)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if resp != nil && len(resp) > 0 {
		return resp, nil
	}

	logMap := make(map[string]interface{})
	defer func() {
		if len(resp) > 0 {
			_ = SetList(lastId, resp)
		}
		c.Info("从DB获取文章列表", logMap, nil)
	}()

	articles, err := new(entity.Article).FindByLastId(lastId, 6)
	if err != nil {
		return nil, err
	}
	logMap["articlesLen"] = len(articles)
	resp = make([]*resp2.GetArticlesResp, 0, len(articles))

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
		data := &resp2.GetArticlesResp{
			Article: v,
		}
		if articleCommentCount, ok := articleCommentCountMap[v.Id]; ok {
			data.CommentCount = articleCommentCount
		}
		resp = append(resp, data)
	}

	return resp, nil
}

func Get(id int) (*resp2.GetArticlesResp, error) {
	resp, err := GetById(id)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if resp != nil && resp.Id > 0 {
		return resp, nil
	}

	article, err := new(entity.Article).GetById(id)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = Set(id, resp)
	}()
	comment, err := new(entity.Comment).GetArticlesWebCommentCount(article.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	resp = &resp2.GetArticlesResp{
		Article:      article,
		CommentCount: comment,
	}

	return resp, nil
}
