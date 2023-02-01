package article

import (
	resp2 "blog-api/controller/resp"
	"blog-api/pkg/goredis"
	"encoding/json"
	"fmt"
	"time"
)

const (
	ListStrKey          = "blog:api:article:list:%d"
	SingleArticleStrKey = "blog:api:article:%d"
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
	return goredis.Client.Set(fmt.Sprintf(ListStrKey, lastId), string(b), time.Hour*24*7).Err()
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
	return goredis.Client.Set(fmt.Sprintf(SingleArticleStrKey, id), string(b), time.Hour*24*7).Err()
}
