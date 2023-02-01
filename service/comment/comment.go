package comment

import (
	"blog-api/controller/resp"
	"blog-api/core"
	"blog-api/entity"
	"blog-api/pkg/util"
	"github.com/go-redis/redis/v7"
)

func FindByArticleIdLastId(c *core.Context, articleId, lastId int) ([]*resp.GetCommentsResp, error) {
	logMap := make(map[string]interface{})
	logMap["articleId"] = articleId
	logMap["lastId"] = lastId
	list, err := GetListByArticleIdLastId(articleId, lastId)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if list != nil {
		return list, nil
	}
	defer func() {
		if len(list) > 0 {
			_ = SetList(articleId, lastId, list)
		}
		c.Info("从DB获取文章列表", logMap, nil)
	}()
	comments, err := new(entity.Comment).FindByArticleIdLastId(articleId, lastId, 15)
	if err != nil {
		c.ErrorL("获取评论失败", logMap, err.Error())
		return nil, err
	}
	list = make([]*resp.GetCommentsResp, 0, len(comments))

	for _, v := range comments {
		emailMd5Str := v.Email
		if v.Email == "" {
			v.Email = "empty"
		}
		data := &resp.GetCommentsResp{
			Name:      v.Name,
			Avatar:    "https://www.gravatar.com/avatar/" + util.EncryptMd5(emailMd5Str),
			Url:       v.Url,
			Type:      v.Type,
			Content:   v.Content,
			ArticleId: v.ArticleId,
			Pid:       v.Pid,
			Ppid:      v.Ppid,
			CreatedAt: v.CreatedAt,
		}
		list = append(list, data)
	}

	return list, nil
}

func GetCommentCount(articleId int) (int64, error) {
	count, err := GetCommentCountByArticleId(articleId)
	if err != nil && err != redis.Nil {
		return 0, err
	}
	if err == nil {
		return count, nil
	}
	defer func() {
		_ = SetCommentCount(articleId, count)
	}()
	count, err = new(entity.Comment).GetArticlesWebCommentCount(articleId)

	return count, err
}
