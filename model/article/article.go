package article

import (
	"fmt"
	"net/http"

	"go-blog/server/errno"

	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {
	var r struct {
		Title string `json:"title"`
		Image string `json:"image"`
		Html string `json:"html"`
		Con string `json:"con"`
	}

	var err error
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": errno.ServerError})
		return
	}

	if r.Title == "" {
		err = errno.New(errno.ServerError, fmt.Errorf("title can not null"))
	}

	code, message := errno.ReturnErr(err)
	c.JSON(http.StatusOK, gin.H{"code": code, "message": message})
}