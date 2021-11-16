package admin

import (
	"github.com/gin-gonic/gin"
	"go-blog/controller"
	"go-blog/model/article"
	"go-blog/model/comment"
	"go-blog/server/errno"
	"strconv"
)

type ArticleRequest struct {
	Title  string `json:"title" binding:"required"`
	Image  string `json:"image"`
	Html   string `json:"html" binding:"required"`
	Con    string `json:"con" binding:"required"`
	Tags   string `json:"tags" binding:"required"`
	Status int    `json:"status"`
	Source int    `json:"source"`
}

func AddArticle(c *gin.Context) {
	var r ArticleRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		controller.Response(c, errno.BindError, err)
		return
	}

	res := &article.Articles{
		Title: r.Title,
		Image: r.Image,
		Html:  r.Html,
		Con:   r.Con,
		Tags:  r.Tags,
	}
	if err := res.Create(); err != nil {
		controller.Response(c, errno.DatabaseError, nil)
		return
	}

	controller.Response(c, nil, res)

	return
}

func DelArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := article.Delete(id); err != nil {
		controller.Response(c, errno.DatabaseError, nil)
		return
	}

	controller.Response(c, nil, nil)

	return
}

func UpdateArticle(c *gin.Context) {
	var r ArticleRequest
	id, _ := strconv.Atoi(c.Param("id"))
	if err := c.ShouldBindJSON(&r); err != nil {
		controller.Response(c, errno.BindError, nil)
		return
	}

	res := &article.Articles{
		Title: r.Title,
		Image: r.Image,
		Html:  r.Html,
		Con:   r.Con,
		Tags:  r.Tags,
	}
	if err := res.Update(id); err != nil {
		controller.Response(c, errno.DatabaseError, nil)
		return
	}

	controller.Response(c, nil, res)

	return
}

func GetArticles(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("last_id", "0"))
	if err != nil {
		controller.Response(c, err, nil)
		return
	}

	articles, err := article.Index(page, 15)

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
