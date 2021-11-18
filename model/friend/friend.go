package friend

import (
	"go-blog/model"
	"time"
)

type Friend struct {
	model.Base
	Name  string `json:"name"`
	Email string `json:"email"`
	Url   string `json:"url"`
}

func (h *Friend) Create() error {
	h.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	h.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	return model.DB.Self.Create(h).Error
}

func Delete(id int) error {
	h := Friend{}
	h.Id = id

	return model.DB.Self.Delete(h).Error
}

func (h *Friend) Update(id int) error {
	h.Id = id
	h.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	return model.DB.Self.Model(h).Update(h).Error
}

func WebIndex() ([]*Friend, error) {
	harems := make([]*Friend, 0)

	err := model.DB.Self.Order("id desc").Find(&harems).Error

	return harems, err
}

func AdminIndex(page, limit int) ([]*Friend, int, error) {
	if limit == 0 {
		limit = 6
	}

	friend := make([]*Friend, 0)
	var count int

	if err := model.DB.Self.Model(&Friend{}).Count(&count).Error; err != nil {
		return friend, count, err
	}
	err := model.DB.Self.Offset((page - 1) * limit).Limit(limit).Order("id desc").Find(&friend).Error

	return friend, count, err
}

func Show(id int) (*Friend, error) {
	harems := &Friend{}
	if err := model.DB.Self.First(&harems, id).Error; err != nil {
		return harems, err
	}

	return harems, nil
}
