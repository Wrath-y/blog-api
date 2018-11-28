package article_dao

import (
	"go-blog/model"
	"go-blog/model/article"
	"go-blog/struct/article-struct"
)



func List(offset, limit int) ([]*article.ArticleModel, uint64, error) {
	if limit == 0 {
		limit = 6
	}

	articles := make([]*article.ArticleModel, 0)
	var count uint64

	if err := model.DB.Self.Model(&article.ArticleModel{}).Count(&count).Error; err != nil {
		return articles, count, err
	}

	if err := model.DB.Self.Offset(offset).Limit(limit).Order("id desc").Find(&articles).Error; err != nil {
		return articles, count, err
	}

	return articles, count, nil
}