package comment

import (
	"blog-api/controller/resp"
	"blog-api/pkg/goredis"
	"encoding/json"
	"fmt"
	"time"
)

const (
	ListStrKey  = "blog:comment:list:%d:%d"
	CountStrKey = "blog:comment:count:%d"
)

func GetListByArticleIdLastId(articleId, lastId int) ([]*resp.GetCommentsResp, error) {
	var res []*resp.GetCommentsResp
	b, err := goredis.Client.Get(fmt.Sprintf(ListStrKey, articleId, lastId)).Bytes()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func SetList(articleId, lastId int, resp []*resp.GetCommentsResp) error {
	b, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	return goredis.Client.Set(fmt.Sprintf(ListStrKey, articleId, lastId), string(b), time.Hour*24*7).Err()
}

func GetCommentCountByArticleId(articleId int) (int64, error) {
	return goredis.Client.Get(fmt.Sprintf(CountStrKey, articleId)).Int64()
}

func SetCommentCount(articleId int, count int64) error {
	return goredis.Client.Set(fmt.Sprintf(CountStrKey, articleId), count, time.Hour*24*7).Err()
}

func ClearCommentCache(articleId, lastId int) error {
	return goredis.Client.Unlink(fmt.Sprintf(CountStrKey, articleId), fmt.Sprintf(ListStrKey, articleId, lastId)).Err()
}
