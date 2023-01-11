package entity

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Base struct {
	Id        int    `json:"id"`
	UpdatedAt string `json:"updated_at"`
	CreatedAt string `json:"created_at"`
}
