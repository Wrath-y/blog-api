package article_dao

import (
	"go-blog/model"
)

type Articles struct {
	Id 	  uint64 `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
	Html  string `json:"html"`
	Con   string `json:"con"`
}

func List(offset, limit int) ([]*Articles, uint64, error) {
	if limit == 0 {
		limit = 6
	}

	articles := make([]*Articles, 0)
	var count uint64

	if err := model.DB.Self.Model(&Articles{}).Count(&count).Error; err != nil {
		return articles, count, err
	}

	if err := model.DB.Self.Offset(offset).Limit(limit).Order("id desc").Find(&articles).Error; err != nil {
		return articles, count, err
	}

	return articles, count, nil
}