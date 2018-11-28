package article_struct

import "sync"

type ArticleInfo struct {
	Id 	  uint64 `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
	Html  string `json:"html"`
	Con   string `json:"con"`
}

type List struct {
	Lock  *sync.Mutex
	IdMap map[uint64]*ArticleInfo
}