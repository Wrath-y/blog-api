package entity

import (
	"blog-api/pkg/db"
)

type FriendLink struct {
	*Base
	Name  string `json:"name"`
	Email string `json:"email"`
	Url   string `json:"url"`
}

func (*FriendLink) TableName() string {
	return "friend_link"
}

func (*FriendLink) FindAll() ([]*FriendLink, error) {
	var harems []*FriendLink
	return harems, db.Orm.Raw("select * from friend_link order by id desc").Find(&harems).Error
}
