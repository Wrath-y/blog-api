package admin

import (
	"github.com/gin-gonic/gin"
	"go-blog/controller"
	"go-blog/model/article"
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
	if err := res.Update(id); err != nil {
		controller.Response(c, errno.DatabaseError, err)
		return
	}

	controller.Response(c, nil, res)

	return
}

func GetArticles(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		controller.Response(c, errno.ServerError, err)
		return
	}

	articles, count, err := article.AdminIndex(page, 6)
	if err != nil {
		controller.Response(c, errno.DatabaseError, err)
		return
	}

	controller.Response(c, nil, map[string]interface{}{
		"list":  articles,
		"count": count,
	})

	return
}

func GetArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := article.Show(id)
	if err != nil {
		controller.Response(c, errno.DatabaseError, err)
		return
	}
	controller.Response(c, nil, res)

	return
}
