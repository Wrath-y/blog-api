package req_harem

import (
	"go-blog/model/harem"
)

type Response struct {
	Count int            `json:"count"`
	Data  []*harem.Harem `json:"data"`
}
