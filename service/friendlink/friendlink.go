package friendlink

import (
	"blog-api/controller/resp"
	"blog-api/entity"
	"github.com/go-redis/redis/v7"
)

func FindAll() ([]*resp.GetFriendLinkResp, error) {
	list, err := GetList()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if list != nil {
		return list, nil
	}

	defer func() {
		if list != nil {
			_ = SetList(list)
		}
	}()
	friendLinks, err := new(entity.FriendLink).FindAll()
	if err != nil {
		return nil, err
	}
	for _, v := range friendLinks {
		list = append(list, &resp.GetFriendLinkResp{
			Name: v.Name,
			Url:  v.Url,
		})
	}

	return list, nil
}
