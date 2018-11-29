package article

import (
	"go-blog/dao/article-dao"
	"go-blog/struct/article-struct"
	"sync"
)

func Index(page, limit int) ([]*article_struct.Article, int, error) {
	data := make([]*article_struct.Article, 0)
	articles, count, err := article_dao.List(page, limit)
	if err != nil {
		return nil, count, err
	}

	ids := []int{}
	for _, article := range articles {
		ids = append(ids, article.Id)
	}

	wg := sync.WaitGroup{}
	articleList := article_struct.List{
		Lock:	new(sync.Mutex),
		IdMap:	make(map[int]*article_struct.Article, len(articles)),
	}

	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	for _, a := range articles {
		wg.Add(1)
		go func(a *article_dao.Articles) {
			defer wg.Done()

			articleList.Lock.Lock()
			defer articleList.Lock.Unlock()
			articleList.IdMap[a.Id] = &article_struct.Article{
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
		data = append(data, articleList.IdMap[id])
	}

	return data, count, nil
}