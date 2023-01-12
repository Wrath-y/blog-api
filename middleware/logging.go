package middleware

import (
	"blog-api/core"
	"blog-api/errcode"
	"blog-api/pkg/def"
	"blog-api/pkg/logging"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

const bodyLimitKB = 5000

type BodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *BodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *BodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Logging(c *core.Context) {
	start := time.Now()

	raw, _ := c.GetRawData()
	c.Request.Body = io.NopCloser(bytes.NewBuffer(raw))

	w := &BodyLogWriter{ResponseWriter: c.Writer, body: bytes.NewBufferString("")}
	c.Writer = w

	logger := logging.New()
	logger.SetRequestID(c.GetString(def.RequestID))
	logger.SetV1(c.GetString("v1"))
	c.Logger = logger

	rawKB := len(raw) / 1024 // => to KB
	if rawKB > bodyLimitKB {
		c.Info("接口请求与响应", string(raw[:1024]), nil)
		c.FailWithErrCode(errcode.WebBodyTooLarge.WithDetail(fmt.Sprintf("消息限制%dKB, 本消息%dKB", bodyLimitKB, rawKB)), nil)
		return
	}

	c.Next()

	logger.SetV2(c.GetString("v2"))
	logger.SetV3(c.GetString("v3"))

	reqBody := new(interface{})
	if err := json.Unmarshal(raw, reqBody); err != nil {
		*reqBody = string(raw)
	}

	request := map[string]interface{}{
		"method":      c.Request.Method,
		"header":      convertHTTPHeader(c.Request.Header),
		"body":        *reqBody,
		"request_uri": c.Request.RequestURI,
	}

	respBody := new(interface{})
	rawRespBody := w.body.Bytes()
	if err := json.Unmarshal(rawRespBody, respBody); err != nil {
		*respBody = string(rawRespBody)
	}

	c.Info("接口请求与响应", request, *respBody, start)
}

func convertHTTPHeader(header http.Header) map[string]interface{} {
	h := make(map[string]interface{})
	for k, v := range header {
		if len(v) > 1 {
			h[k] = v
		} else {
			h[k] = v[0]
		}
	}
	return h
}
