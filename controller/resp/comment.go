package resp

import "time"

type GetCommentsResp struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	Url       string    `json:"url"`
	Content   string    `json:"content"`
	ArticleId int       `json:"article_id"`
	Pid       int       `json:"pid"`
	CreatedAt time.Time `json:"created_at"`
}
