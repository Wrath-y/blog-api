package commentStruct

import "go-blog/model/comment"

type Response struct {
	Count int	 `json:"count"`
	Data  []*comment.Comment `json:"data"`
}