package article

import (
	"go-blog/dao/article-dao"
	"go-blog/struct/article-struct"
	"sync"
)

type ArticleModel struct {
	Id 	  uint64 `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
	Html  string `json:"html"`
	Con   string `json:"con"`
}

func Index(offset, limit int) ([]*article_struct.ArticleInfo, uint64, error) {
	infos := make([]*article_struct.ArticleInfo, 0)
	articles, count, err := article_dao.List(offset, limit)
	if err != nil {
		return nil, count, err
	}

	ids := []uint64{}
	for _, article := range articles {
		ids = append(ids, article.Id)
	}

	wg := sync.WaitGroup{}
	articleList := article_struct.List{
		Lock:	new(sync.Mutex),
		IdMap:	make(map[uint64]*article_struct.ArticleInfo, len(articles)),
	}

	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	for _, a := range articles {
		wg.Add(1)
		go func(a *ArticleModel) {
			defer wg.Done()

			articleList.Lock.Lock()
			defer articleList.Lock.Unlock()
			articleList.IdMap[a.Id] = &article_struct.ArticleInfo{
				Id:		a.Id,
				Title:	a.Title,
				Image:  a.Image,
				Html:   a.Html,
				Con:	a.Con,
			}
		}(a)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		return nil, count, err
	}

	for _, id := range ids {
		infos = append(infos, articleList.IdMap[id])
	}

	return infos, count, nil
}