package req_spider

import (
	"github.com/gin-gonic/gin"
	"go-blog/server/errno"
	"reflect"
)

type UpdateImgRequest struct {
	Cookie string `json:"cookie"`
}

func (r UpdateImgRequest) Validate(_ *gin.Context) error {
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
