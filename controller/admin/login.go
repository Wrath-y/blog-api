package admin

import (
	"github.com/gin-gonic/gin"
	"go-blog/controller"
	"go-blog/model/administrator"
	"go-blog/server/auth"
	"go-blog/server/errno"
	"go-blog/server/token"
)

type UserRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var r UserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		controller.Response(c, errno.BindError, nil)

		return
	}

	// Get the administrator information by the login username.
	a, err := administrator.GetUserByName(r.Account)
	if err != nil {
		controller.Response(c, errno.UserNotFoundErr, nil)

		return
	}

	if err != nil {
		controller.Response(c, errno.ServerError, nil)

		return
	}

	// Compare the login password with the administrator password.
	if err := auth.Compare(a.Password, r.Password); err != nil {
		controller.Response(c, errno.PasswordIncorrectErr, nil)
	}

	// Sign the json web token.
	t, err := token.Sign(c, token.Context{ID: a.Id, Account: a.Account}, "")
	if err != nil {
		controller.Response(c, errno.TokenErr, nil)

		return
	}

	controller.Response(c, nil, administrator.Token{Token: t})

	return
}
