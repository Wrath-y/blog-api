package middleware

import (
	"blog-api/core"
	"time"
)

func TimeLocation(c *core.Context) {
	l, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		l = time.FixedZone("CST", 8*3600)
	}

	c.TimeLocation = l
}
