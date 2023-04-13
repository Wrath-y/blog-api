package article

import (
	"blog-api/controller/resp"
	"blog-api/entity"
	"blog-api/service/comment"
	"github.com/go-redis/redis/v7"
)

func GetBaseInfo(id int) (*resp.GetArticleBaseInfoResp, error) {
	res, err := GetBaseInfoByID(id)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if res != nil {
		return res, nil
	}

	article, err := new(entity.Article).GetBaseInfoById(id)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = SetBaseInfo(id, res)
	}()

	commentCount, _ := comment.GetCommentCount(id)

	res = &resp.GetArticleBaseInfoResp{
		Hits:         article.Hits,
		CommentCount: commentCount,
	}

	return res, nil
}
