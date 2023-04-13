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
	Hits   int64  `json:"hits"`
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
	if lastId == 0 {
		return article, db.Orm.Raw("select * from article where status = 1 order by id desc limit ?", limit).Find(&article).Error
	}
	if lastId == -1 {
		return article, db.Orm.Raw("select * from article where status = 1 order by id desc").Find(&article).Error
	}
	return article, db.Orm.Raw("select * from article where id < ? and status = 1 order by id desc limit ?", lastId, limit).Find(&article).Error
}

func (*Article) FindAll() ([]*Article, error) {
	var article []*Article
	return article, db.Orm.Raw("select * from article order by id desc").Find(&article).Error
}

func (*Article) GetById(id int) (*Article, error) {
	article := new(Article)
	return article, db.Orm.First(&article, id).Error
}

func (*Article) GetBaseInfoById(id int) (*Article, error) {
	article := new(Article)
	return article, db.Orm.Raw("select id, hits from article where id = ?", id).First(&article).Error
}

func (*Article) HitsIncr(id int) error {
	return db.Orm.Exec("update article set hits = hits + 1 where id = ?", id).Error
}
