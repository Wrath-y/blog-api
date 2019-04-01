package article

import (
	"go-blog/model"
	"time"
)

type Articles struct {
	model.Base
	Title 		string  `json:"title"`
	Image 		string  `json:"image"`
	Html  		string  `json:"html"`
	Con   		string  `json:"con"`
	Tag   		string  `json:"tag"`
	Status		int     `json:"status"`
	Source		int     `json:"source"`
}

func (a *Articles) Create() error {
	a.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	a.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	return model.DB.Self.Create(a).Error
}

func Delete(id int) error {
	a := Articles{}
	a.Id = id

	return model.DB.Self.Delete(a).Error
}

func (a *Articles) Update(id int) error {
	a.Id = id
	a.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	return model.DB.Self.Model(a).Update(a).Error
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

func Show(id int) (*Articles, error) {
	articles := &Articles{}
	if err := model.DB.Self.First(&articles, id).Error; err != nil {
		return articles, err
	}

	return articles, nil
}