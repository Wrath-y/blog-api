package administrator

import "go-blog/model"

type Administrators struct {
	model.Base
	Account 		string  `json:"account"`
	Password 		string  `json:"password"`
}