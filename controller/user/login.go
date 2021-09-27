package user

import (
	"github.com/gin-gonic/gin"
	"go-blog/model"
	"go-blog/model/administrator"
	"go-blog/req_struct"
	"go-blog/req_struct/req_login"
	"go-blog/server/auth"
	"go-blog/server/errno"
	"go-blog/server/token"
)

func Login(c *gin.Context) {
	var r req_login.Request
	if err := c.Bind(&r); err != nil {
		req_struct.Response(c, errno.BindError, nil)

		return
	}

	// Get the administrator information by the login username.
	a, err := administrator.GetUserByName(r.Account)
	if err != nil {
		req_struct.Response(c, errno.UserNotFoundErr, nil)

		return
	}

	if err != nil {
		req_struct.Response(c, errno.ServerError, nil)

		return
	}

	// Compare the login password with the administrator password.
	if err := auth.Compare(a.Password, r.Password); err != nil {
		req_struct.Response(c, errno.PasswordIncorrectErr, nil)
	}

	// Sign the json web token.
	t, err := token.Sign(c, token.Context{ID: a.Id, Account: a.Account}, "")
	if err != nil {
		req_struct.Response(c, errno.TokenErr, nil)

		return
	}

	req_struct.Response(c, nil, model.Token{Token: t})

	return
}
