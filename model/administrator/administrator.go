package administrator

import (
	"go-blog/model"
)

type Administrators struct {
	model.Base
	Account  string `json:"account"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

func GetUserByName(account string) (*Administrators, error) {
	a := &Administrators{}

	return a, model.DB.Self.Where("account = ?", account).First(&a).Error
}
