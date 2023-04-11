package articleseo

import (
	resp2 "blog-api/controller/resp"
	"blog-api/pkg/goredis"
	"encoding/json"
	"fmt"
	"time"
)

const (
	ListStrKey = "blog:article-seo:%d"
)

func GetListById(id int) ([]*resp2.GetArticleSeoResp, error) {
	var resp []*resp2.GetArticleSeoResp
	b, err := goredis.Client.Get(fmt.Sprintf(ListStrKey, id)).Bytes()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func SetList(id int, resp []*resp2.GetArticleSeoResp) error {
	b, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	return goredis.Client.Set(fmt.Sprintf(ListStrKey, id), string(b), time.Hour*24*7).Err()
}
