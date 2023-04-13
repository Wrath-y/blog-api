package article

import (
	"blog-api/controller/resp"
	"blog-api/core"
	"blog-api/entity"
	"blog-api/service/comment"
	"github.com/go-redis/redis/v7"
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

	commentCounts := make(map[int]int64)
	for _, v := range articleIds {
		commentCount, err := comment.GetCommentCount(v)
		if err != nil {
			logMap["cur_article_id"] = v
			c.ErrorL("获取评论失败", logMap, err.Error())
			commentCount = 0
		}
		commentCounts[v] = commentCount
	}

	articleCommentCountMap := make(map[int]int64)
	for articleId, v := range commentCounts {
		articleCommentCountMap[articleId] = v
	}

	for _, v := range articles {
		data := &resp.GetArticlesResp{
			ID:        v.Id,
			Title:     v.Title,
			Image:     v.Image,
			Html:      v.Html,
			Tags:      v.Tags,
			Hits:      v.Hits,
			CreatedAt: v.CreatedAt,
		}
		if articleCommentCount, ok := articleCommentCountMap[v.Id]; ok {
			data.CommentCount = articleCommentCount
		}
		list = append(list, data)
	}

	return list, nil
}

func All(c *core.Context) ([]*resp.GetArticlesResp, error) {
	list := make([]*resp.GetArticlesResp, 0)
	logMap := make(map[string]interface{})

	articles, err := new(entity.Article).FindAll()
	if err != nil {
		return nil, err
	}
	logMap["articlesLen"] = len(articles)

	articleIds := make([]int, 0, len(articles))
	for _, v := range articles {
		articleIds = append(articleIds, v.Id)
	}
	logMap["articleIds"] = articleIds

	commentCounts := make(map[int]int64)
	for _, v := range articleIds {
		commentCount, err := comment.GetCommentCount(v)
		if err != nil {
			logMap["cur_article_id"] = v
			c.ErrorL("获取评论失败", logMap, err.Error())
			commentCount = 0
		}
		commentCounts[v] = commentCount
	}

	articleCommentCountMap := make(map[int]int64)
	for articleId, v := range commentCounts {
		articleCommentCountMap[articleId] = v
	}

	for _, v := range articles {
		data := &resp.GetArticlesResp{
			ID:        v.Id,
			Title:     v.Title,
			Image:     v.Image,
			Html:      v.Html,
			Tags:      v.Tags,
			Hits:      v.Hits,
			CreatedAt: v.CreatedAt,
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
	if res != nil && res.ID > 0 {
		return res, nil
	}

	article, err := new(entity.Article).GetById(id)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = Set(id, res)
	}()

	commentCount, _ := comment.GetCommentCount(id)

	res = &resp.GetArticlesResp{
		ID:           article.Id,
		Title:        article.Title,
		Image:        article.Image,
		Html:         article.Html,
		Tags:         article.Tags,
		Hits:         article.Hits,
		CreatedAt:    article.CreatedAt,
		CommentCount: commentCount,
	}

	return res, nil
}

func HitsIncr(id int) error {
	if err := new(entity.Article).HitsIncr(id); err != nil {
		return err
	}
	if err := CacheHitsIncr(id); err != nil {
		return err
	}
	return nil
}
