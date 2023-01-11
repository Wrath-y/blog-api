package comment

import (
	"go-blog/model"
	"go-blog/pkg/db"
	"time"
)

type Comment struct {
	model.Base
	Name      string `json:"name" gorm:"name"`
	Email     string `json:"email" gorm:"email"`
	Url       string `json:"url" gorm:"url"`
	Type      int    `json:"type" gorm:"type"`
	Content   string `json:"content" gorm:"content"`
	ArticleId int    `json:"article_id" gorm:"article_id"`
	Pid       int    `json:"pid" gorm:"pid"`
	Ppid      int    `json:"ppid" gorm:"ppid"`
}

type ArticlesWebCommentCount struct {
	ArticleId    int `json:"article_id"`
	CommentCount int `json:"comment_count"`
}

func AdminIndex(page, limit int) ([]*Comment, int64, error) {
	if limit == 0 {
		limit = 6
	}

	var count int64
	if err := db.Orm.Model(&Comment{}).Count(&count).Error; err != nil {
		return nil, count, err
	}

	comments := make([]*Comment, 0)
	err := db.Orm.Offset((page - 1) * limit).Limit(limit).Find(&comments).Error

	return comments, count, err
}

func IndexBuyArticleId(articleId, lastId, limit int) ([]*Comment, error) {
	if limit == 0 {
		limit = 6
	}

	comments := make([]*Comment, 0)
	err := db.Orm.Where("id > ?", lastId).Where("article_id = ?", articleId).Limit(limit).Find(&comments).Error

	return comments, err
}

func Delete(id int) error {
	c := Comment{}
	c.Id = id

	return db.Orm.Delete(c).Error
}

func (c *Comment) Create() error {
	c.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	c.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	return db.Orm.Create(c).Error
}

func GetArticlesWebCommentCounts(articleIds []int) ([]*ArticlesWebCommentCount, error) {
	articlesWebCommentCount := make([]*ArticlesWebCommentCount, 0)
	err := db.Orm.Table("comments").Select("article_id, COUNT(id) as comment_count").Where("article_id in (?)", articleIds).Group("article_id").Find(&articlesWebCommentCount).Error

	return articlesWebCommentCount, err
}

func GetArticlesWebCommentCount(articleId int) (*ArticlesWebCommentCount, error) {
	articlesWebCommentCount := new(ArticlesWebCommentCount)
	err := db.Orm.Table("comments").Select("article_id, COUNT(id) as comment_count").Where("article_id = (?)", articleId).Group("article_id").Find(&articlesWebCommentCount).Error

	return articlesWebCommentCount, err
}
