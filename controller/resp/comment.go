package resp

import "time"

type GetCommentsResp struct {
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	Url       string    `json:"url"`
	Type      int       `json:"type"`
	Content   string    `json:"content"`
	ArticleId int       `json:"article_id"`
	Pid       int       `json:"pid"`
	Ppid      int       `json:"ppid"`
	CreatedAt time.Time `json:"created_at"`
}
