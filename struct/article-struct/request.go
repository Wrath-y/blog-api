package article_struct

type Request struct {
	Title string `json:"title"`
	Image string `json:"image"`
	Html string `json:"html"`
	Con string `json:"con"`
}