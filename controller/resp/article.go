package resp

import "blog-api/entity"

type GetArticlesResp struct {
	Article      *entity.Article `json:"article"`
	CommentCount int             `json:"comment_count"`
}
