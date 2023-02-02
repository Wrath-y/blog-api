package entity

import (
	"blog-api/pkg/db"
)

type Comment struct {
	*Base
	Name      string `json:"name" gorm:"name"`
	Email     string `json:"email" gorm:"email"`
	Url       string `json:"url" gorm:"url"`
	Content   string `json:"content" gorm:"content"`
	ArticleId int    `json:"article_id" gorm:"article_id"`
	Pid       int    `json:"pid" gorm:"pid"`
}

type ArticlesWebCommentCount struct {
	ArticleId    int   `json:"article_id"`
	CommentCount int64 `json:"comment_count"`
}

func (*Comment) TableName() string {
	return "comment"
}

func (c *Comment) Create() error {
	return db.Orm.Create(c).Error
}

func (*Comment) FindByArticleIdLastId(articleId, lastId, limit int) ([]*Comment, error) {
	if limit == 0 {
		limit = 6
	}
	var comments []*Comment
	if lastId == 0 {
		return comments, db.Orm.Raw("select * from comment where article_id = ? order by id desc limit ?", articleId, limit).Find(&comments).Error
	}
	return comments, db.Orm.Raw("select * from comment where id < ? and article_id = ? order by id desc limit ?", lastId, articleId, limit).Find(&comments).Error
}

func (*Comment) GetArticlesWebCommentCounts(articleIds []int) ([]*ArticlesWebCommentCount, error) {
	var articlesWebCommentCount []*ArticlesWebCommentCount
	return articlesWebCommentCount, db.Orm.Raw("select article_id, COUNT(id) as comment_count from comment where article_id in ? group by article_id", articleIds).Find(&articlesWebCommentCount).Error
}

func (*Comment) GetArticlesWebCommentCount(articleId int) (int64, error) {
	var count int64
	return count, db.Orm.Raw("select COUNT(id) as comment_count from comment where article_id = ?", articleId).Count(&count).Error
}
