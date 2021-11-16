package admin

import (
	"github.com/gin-gonic/gin"
	"go-blog/controller"
	"go-blog/model/harem"
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

	res := &harem.Harem{
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

	if err := harem.Delete(id); err != nil {
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

	res := &harem.Harem{
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
	data, err := harem.Index(page, 15)

	if err != nil {
		controller.Response(c, err, nil)
		return
	}

	controller.Response(c, nil, data)

	return
}

func GetFriend(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := harem.Show(id)
	if err != nil {
		controller.Response(c, errno.DatabaseError, nil)
		return
	}
	controller.Response(c, nil, res)

	return
}
