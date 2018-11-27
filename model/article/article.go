package article

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-blog/server/errno"
	"go-blog/struct/articleStruct"
)

func Create(c *gin.Context) {
	var r articleStruct.Request
	var err error
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": errno.BindError})

		return
	}

	if r.Title == "" {
		err = errno.New(errno.TitleError, fmt.Errorf("title can not null"))
	}

	code, message := errno.ReturnErr(err)
	c.JSON(http.StatusOK, gin.H{"code": code, "message": message})
}