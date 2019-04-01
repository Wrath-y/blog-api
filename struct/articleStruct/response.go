package articleStruct

import (
	"go-blog/model/article"
)

type Response struct {
	Count int	 `json:"count"`
	Data  []*article.Articles `json:"data"`
}