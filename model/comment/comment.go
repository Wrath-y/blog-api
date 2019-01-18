package comment

import "go-blog/model"

type Comment struct {
	model.Base
	Name	string `json:"name"`
	Email	string `json:"email"`
	Url		string `json:"url"`
	Type    string `json:"type"`
	Content string `json:"content"`
	ArticleId int  `json:"article_id"`
	Pid     int	   `json:"pid"`
	Ppid    int    `json:"ppid"`
}

func Index(page, limit int) ([]*Comment, int, error) {
	if limit == 0 {
		limit = 6
	}

	comments := make([]*Comment, 0)
	var count int

	if err := model.DB.Self.Model(&Comment{}).Count(&count).Error; err != nil {
		return comments, count, err
	}

	if err := model.DB.Self.Offset((page - 1) * limit).Limit(limit).Find(&comments).Error; err != nil {
		return comments, count ,err
	}

	return comments, count, nil
}

func Delete(id int) error {
	c := Comment{}
	c.Id = id

	return model.DB.Self.Delete(c).Error
}