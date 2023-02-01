package resp

import "blog-api/entity"

type GetArticlesResp struct {
	*entity.Article
	CommentCount int64 `json:"comment_count"`
}
