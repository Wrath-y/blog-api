package article

import (
	resp2 "blog-api/controller/resp"
	"blog-api/pkg/goredis"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const (
	ListStrKey                  = "blog:article:list:%d"
	SingleArticleStrKey         = "blog:article:%d"
	SingleArticleBaseInfoStrKey = "blog:article:baseinfo:%d"
)

func GetListByLastId(lastId int) ([]*resp2.GetArticlesResp, error) {
	var resp []*resp2.GetArticlesResp
	b, err := goredis.Client.Get(fmt.Sprintf(ListStrKey, lastId)).Bytes()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func SetList(lastId int, resp []*resp2.GetArticlesResp) error {
	b, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	return goredis.Client.Set(fmt.Sprintf(ListStrKey, lastId), string(b), time.Hour*24*180).Err()
}

func GetById(id int) (*resp2.GetArticlesResp, error) {
	var resp *resp2.GetArticlesResp
	b, err := goredis.Client.Get(fmt.Sprintf(SingleArticleStrKey, id)).Bytes()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func Set(id int, resp *resp2.GetArticlesResp) error {
	b, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	return goredis.Client.Set(fmt.Sprintf(SingleArticleStrKey, id), string(b), time.Hour*24*180).Err()
}

func GetBaseInfoByID(id int) (*resp2.GetArticleBaseInfoResp, error) {
	resp := new(resp2.GetArticleBaseInfoResp)
	m, err := goredis.Client.HGetAll(fmt.Sprintf(SingleArticleBaseInfoStrKey, id)).Result()
	if err != nil {
		return nil, err
	}
	if hits, ok := m["hits"]; ok {
		h, _ := strconv.Atoi(hits)
		resp.Hits = int64(h)
	}
	if cc, ok := m["comment_count"]; ok {
		c, _ := strconv.Atoi(cc)
		resp.CommentCount = int64(c)
	}

	return resp, nil
}

func SetBaseInfo(id int, resp *resp2.GetArticleBaseInfoResp) error {
	params := make([]interface{}, 0)
	params = append(params, "hits", resp.Hits, "comment_count", resp.CommentCount)

	return goredis.Client.HMSet(fmt.Sprintf(SingleArticleBaseInfoStrKey, id), params...).Err()
}

func CacheHitsIncr(id int) error {
	return goredis.Client.HIncrBy(fmt.Sprintf(SingleArticleBaseInfoStrKey, id), "hits", 1).Err()
}

func CacheCommentCountIncr(id int) error {
	return goredis.Client.HIncrBy(fmt.Sprintf(SingleArticleBaseInfoStrKey, id), "comment_count", 1).Err()
}
