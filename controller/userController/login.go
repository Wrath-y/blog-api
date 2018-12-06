package userController

import (
	"github.com/gin-gonic/gin"
	"go-blog/model/administrator"
	"go-blog/server/auth"
	"go-blog/server/errno"
	"go-blog/server/token"
	"go-blog/struct"
	"go-blog/struct/login-struct"
)

func Login(c *gin.Context) {
	var r login_struct.Request
	if err := c.Bind(&r); err != nil {
		_struct.Response(c, errno.BindError, nil)

		return
	}

	// Get the administrator information by the login username.
	a, err := administrator.GetUserByName(r.Account)
	if err != nil {
		_struct.Response(c, errno.ErrUserNotFound, nil)

		return
	}

	// Compare the login password with the administrator password.
	if err := auth.Compare(a.Password, r.Password); err != nil {
		_struct.Response(c, errno.ErrPasswordIncorrect, nil)
		return
	}

	// Sign the json web token.
	t, err := token.Sign(c, token.Context{ID: a.Id, Username: a.Username}, "")
	if err != nil {
		_struct.Response(c, errno.ErrToken, nil)
		return
	}

	SendResponse(c, nil, model.Token{Token: t})
}