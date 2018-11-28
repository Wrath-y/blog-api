package article_struct

type Response struct {
	Title string `json:"title"`
	Image string `json:"image"`
	Html  string `json:"html"`
	Con   string `json:"con"`
}