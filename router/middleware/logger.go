package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/xjh22222228/ip"
	_struct "go-blog/struct"
	"io/ioutil"
	"strings"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Logger(c *gin.Context) {
	bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = bodyLogWriter

	data := ""
	if c.Request.Method == "POST" {
		body := make([]byte, 1024)
		n, _ := c.Request.Body.Read(body)
		data = strings.Replace(string(body[0:n]), "\r\n    ", "", -1)
		data = strings.Replace(data, "\r\n", "", -1)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body[0:n]))
	}
	if c.Request.Method == "GET" {
		data = c.Request.URL.Query().Encode()
	}

	realIP, _ := ip.V4()
	if realIP == "" {
		realIP, _ = ip.V6()
	}

	log.Info().
		Str("uri", c.Request.RequestURI).
		Str("method", c.Request.Method).
		Str("req_ip", realIP).
		Msg(data)

	//处理请求
	c.Next()

	responseBody := bodyLogWriter.body.String()

	response := _struct.Rule{}
	if responseBody != "" {
		json.Unmarshal([]byte(responseBody), &response)
	}
	responseJson, _ := json.Marshal(response)

	log.Info().
		Str("uri", c.Request.RequestURI).
		Str("method", c.Request.Method).
		Str("req_ip", realIP).
		Msg(string(responseJson))
}