package article

import (
	"go-blog/model"
	"time"
)

type Articles struct {
	model.Base
	Title  string `json:"title"`
	Image  string `json:"image"`
	Html   string `json:"html"`
	Con    string `json:"con"`
	Tags   string `json:"tags"`
	Hits   string `json:"hits"`
	Status int    `json:"status"`
	Source int    `json:"source"`
}

type ArticlesWebIndex struct {
	Articles
	CommentCount int `json:"comment_count"`
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

func Index(lastId, limit int) ([]*ArticlesWebIndex, error) {
	if limit == 0 {
		limit = 6
	}

	articlesWebIndex := make([]*ArticlesWebIndex, 0)
	err := model.DB.Self.Table("articles").Where("id > ?", lastId).Limit(limit).Order("id desc").Find(&articlesWebIndex).Error

	return articlesWebIndex, err
}

func Show(id int) (*Articles, error) {
	articles := &Articles{}
	if err := model.DB.Self.First(&articles, id).Error; err != nil {
		return articles, err
	}

	return articles, nil
}
