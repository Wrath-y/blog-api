package entity

import (
	"blog-api/pkg/db"
)

type Article struct {
	*Base
	Title  string `json:"title"`
	Image  string `json:"image"`
	Html   string `json:"html"`
	Con    string `json:"con"`
	Tags   string `json:"tags"`
	Hits   string `json:"hits"`
	Status int    `json:"status"`
	Source int    `json:"source"`
}

func (*Article) TableName() string {
	return "article"
}

type ArticlesWebIndex struct {
	Article
	CommentCount int `json:"comment_count"`
}

func (*Article) FindByLastId(lastId, limit int) ([]*Article, error) {
	var article []*Article
	return article, db.Orm.Raw("select * from article where id < ? and status = 1 limit ?", lastId, limit).Find(&article).Error
}

func (*Article) GetById(id int) (*Article, error) {
	article := new(Article)
	return article, db.Orm.First(&article, id).Error
}
