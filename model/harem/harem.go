package harem

import (
	"go-blog/model"
	"time"
)

type Harem struct {
	model.Base
	Name  string `json:"name"`
	Email string `json:"email"`
	Url   string `json:"url"`
}

func (h *Harem) Create() error {
	h.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	h.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	return model.DB.Self.Create(h).Error
}

func Delete(id int) error {
	h := Harem{}
	h.Id = id

	return model.DB.Self.Delete(h).Error
}

func (h *Harem) Update(id int) error {
	h.Id = id
	h.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	return model.DB.Self.Model(h).Update(h).Error
}

func Index(page, limit int) ([]*Harem, error) {
	if limit == 0 {
		limit = 6
	}

	harems := make([]*Harem, 0)

	err := model.DB.Self.Offset((page - 1) * limit).Limit(limit).Order("id desc").Find(&harems).Error

	return harems, err
}

func Show(id int) (*Harem, error) {
	harems := &Harem{}
	if err := model.DB.Self.First(&harems, id).Error; err != nil {
		return harems, err
	}

	return harems, nil
}
