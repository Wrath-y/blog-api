package req_article

import (
	"github.com/gin-gonic/gin"
	"go-blog/server/errno"
	"reflect"
)

type Request struct {
	Title  string `json:"title"`
	Image  string `json:"image"`
	Html   string `json:"html"`
	Con    string `json:"con"`
	Tags   string `json:"tags"`
	Status int    `json:"status"`
	Source int    `json:"source"`
}

func (r Request) Validate(_ *gin.Context) error {
	t := reflect.TypeOf(r)
	v := reflect.ValueOf(r)
	for k := 0; k < t.NumField(); k++ {
		switch t.Field(k).Type.String() {
		case "string":
			if v.Field(k).String() == "" {
				err := errno.New(errno.RequestError, " "+t.Field(k).Name+" can not be null")

				return err
			}
		}
	}

	return nil
}
