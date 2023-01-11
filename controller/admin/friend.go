package admin

import (
	"github.com/gin-gonic/gin"
	"go-blog/controller"
	"go-blog/entity/friend"
	"go-blog/server/errno"
	"strconv"
)

type FriendRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Url   string `json:"url" binding:"required"`
}

func AddFriend(c *gin.Context) {
	var r FriendRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		controller.Response(c, errno.BindError, err)
		return
	}

	res := &friend.Friend{
		Name:  r.Name,
		Email: r.Email,
		Url:   r.Url,
	}
	if err := res.Create(); err != nil {
		controller.Response(c, errno.DatabaseError, err)
		return
	}

	controller.Response(c, nil, res)

	return
}

func DelFriend(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := friend.Delete(id); err != nil {
		controller.Response(c, errno.DatabaseError, nil)
		return
	}

	controller.Response(c, nil, nil)

	return
}

func UpdateFriend(c *gin.Context) {
	var r FriendRequest
	id, _ := strconv.Atoi(c.Param("id"))
	if err := c.ShouldBindJSON(&r); err != nil {
		controller.Response(c, errno.BindError, nil)
		return
	}

	res := &friend.Friend{
		Name:  r.Name,
		Email: r.Email,
		Url:   r.Url,
	}
	if err := res.Update(id); err != nil {
		controller.Response(c, errno.DatabaseError, nil)
		return
	}

	controller.Response(c, nil, res)

	return
}

func GetFriends(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		panic(err)
	}

	data, count, err := friend.AdminIndex(page, 15)
	if err != nil {
		controller.Response(c, errno.DatabaseError, err)
		return
	}

	controller.Response(c, nil, map[string]interface{}{
		"list":  data,
		"count": count,
	})

	return
}

func GetFriend(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := friend.Show(id)
	if err != nil {
		controller.Response(c, errno.DatabaseError, nil)
		return
	}
	controller.Response(c, nil, res)

	return
}
