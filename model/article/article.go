package article

import (
	"go-blog/model"
)

type Articles struct {
	Id 	  int `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
	Html  string `json:"html"`
	Con   string `json:"con"`
	Tag   string `json:"tag"'`
}

func (a *Articles) Create() error {
	return model.DB.Self.Create(a).Error
}

func Index(page, limit int) ([]*Articles, int, error) {
	if limit == 0 {
		limit = 6
	}

	articles := make([]*Articles, 0)
	var count int

	if err := model.DB.Self.Model(&Articles{}).Count(&count).Error; err != nil {
		return articles, count, err
	}

	if err := model.DB.Self.Offset((page - 1) * limit).Limit(limit).Order("id desc").Find(&articles).Error; err != nil {
		return articles, count, err
	}

	return articles, count, nil
}