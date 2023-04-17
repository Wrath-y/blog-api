package resp

import "time"

type GetArticlesResp struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Image        string    `json:"image"`
	Html         string    `json:"html"`
	Intro        string    `json:"intro"`
	Tags         string    `json:"tags"`
	Hits         int64     `json:"hits"`
	CreatedAt    time.Time `json:"created_at"`
	CommentCount int64     `json:"comment_count"`
}

type GetArticleBaseInfoResp struct {
	Hits         int64 `json:"hits"`
	CommentCount int64 `json:"comment_count"`
}
