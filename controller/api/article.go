package api

import (
	"github.com/gin-gonic/gin"
	"go-blog/controller"
	"go-blog/model/article"
	"go-blog/model/comment"
	"go-blog/server/errno"
	"strconv"
)

func GetArticles(c *gin.Context) {
	lastId, err := strconv.Atoi(c.DefaultQuery("last_id", "0"))
	if err != nil {
		controller.Response(c, err, nil)
		return
	}

	articles, err := article.WebIndex(lastId, 6)

	if err != nil {
		controller.Response(c, err, nil)
		return
	}

	articleIds := make([]int, 0, len(articles))
	for _, v := range articles {
		articleIds = append(articleIds, v.Id)
	}

	commentCounts, err := comment.GetArticlesWebCommentCounts(articleIds)
	if err != nil {
		controller.Response(c, err, nil)
		return
	}

	articleCommentCountMap := make(map[int]int)
	for _, v := range commentCounts {
		articleCommentCountMap[v.ArticleId] = v.CommentCount
	}

	for _, v := range articles {
		if _, ok := articleCommentCountMap[v.Id]; !ok {
			continue
		}
		v.CommentCount = articleCommentCountMap[v.Id]
	}

	controller.Response(c, nil, articles)

	return
}

func GetArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := article.Show(id)
	if err != nil {
		controller.Response(c, errno.DatabaseError, nil)
		return
	}
	controller.Response(c, nil, res)

	return
}
