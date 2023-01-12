package entity

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type Base struct {
	Id        int       `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
