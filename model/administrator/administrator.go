package administrator

import (
	"go-blog/model"
)

func GetUserByName(account string) (*Administrators, error) {
 	a := &Administrators{}

 	return a, model.DB.Self.Where("account = ?", account).Error
}