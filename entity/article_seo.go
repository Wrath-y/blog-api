package entity

import (
	"blog-api/pkg/db"
)

type ArticleSeo struct {
	*Base
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (*ArticleSeo) TableName() string {
	return "article_seo"
}

func (*ArticleSeo) FindByArticleId(id int) ([]*ArticleSeo, error) {
	var articleSeo []*ArticleSeo
	return articleSeo, db.Orm.Raw("select * from article_seo where article_id = ?", id).Find(&articleSeo).Error
}
