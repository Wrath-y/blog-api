package article_struct

import (
	"sync"
)

type Article struct {
	Id 	  int `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
	Html  string `json:"html"`
	Con   string `json:"con"`
}

type List struct {
	Lock  *sync.Mutex
	IdMap map[int]*Article
}

type Response struct {
	Count int	 `json:"count"`
	Data  []*Article `json:"data"`
}