package friendlink

import (
	"blog-api/controller/resp"
	"blog-api/pkg/goredis"
	"encoding/json"
	"time"
)

const (
	ListStrKey = "blog:friend-link:list"
)

func GetList() ([]*resp.GetFriendLinkResp, error) {
	var res []*resp.GetFriendLinkResp
	b, err := goredis.Client.Get(ListStrKey).Bytes()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func SetList(list []*resp.GetFriendLinkResp) error {
	b, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return goredis.Client.Set(ListStrKey, string(b), time.Hour*24*7).Err()
}
