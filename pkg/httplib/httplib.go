package httplib

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	logTopic string
	logger   Logger
	mu       sync.Mutex
)

const (
	requestID = "Request-Id"
	v1        = "v1"
	v2        = "v2"
	v3        = "v3"
)

type Logger func([]byte)
type Option func(*Client)

type Context interface {
	Set(string, any)
	Get(string) (any, bool)
}

type Client struct {
	Transport *Transport
	Timeout   time.Duration
	Context   Context
}

type Transport struct {
	http.RoundTripper
	NoLog bool
}

// NewClient 使用指定的options返回httplib.Client
func NewClient(options ...Option) *Client {
	c := &Client{
		Transport: &Transport{RoundTripper: http.DefaultTransport},
		Timeout:   3 * time.Second,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

// WithTimeout 指定http请求的超时时间, 默认3秒
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.Timeout = timeout
	}
}

// WithTransport 指定http请求的transport, 默认为http.DefaultTransport
func WithTransport(transport http.RoundTripper) Option {
	return func(c *Client) {
		c.Transport.RoundTripper = transport
	}
}

// WithNoLog 是否调用logger记录请求与响应报文
func WithNoLog(noLog bool) Option {
	return func(c *Client) {
		c.Transport.NoLog = noLog
	}
}

// WithContext 是否调用logger记录请求与响应报文
func WithContext(ctx Context) Option {
	return func(c *Client) {
		c.Context = ctx
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.DoWithContext(req)
}

func (c *Client) DoWithContext(req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Transport: c.Transport,
		Timeout:   c.Timeout,
	}
	return client.Do(req)
}

func Setup(t string, l Logger) {
	mu.Lock()
	defer mu.Unlock()
	if logger != nil {
		return
	}

	logTopic = t
	logger = l
}

func (c *Client) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if logger == nil || c.Transport.NoLog {
		resp, err = c.Transport.RoundTrip(req)
		return
	}

	var rawReqBody []byte
	var rawRespBody []byte

	if req.Body != nil {
		rawReqBody, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewBuffer(rawReqBody))
	}

	start := time.Now()
	resp, err = c.Transport.RoundTrip(req)

	defer func() {
		reqBody := new(interface{})
		if err := json.Unmarshal(rawReqBody, reqBody); err != nil {
			*reqBody = string(rawReqBody)
		}
		request, _ := json.Marshal(map[string]interface{}{
			"method":      req.Method,
			"header":      convertHTTPHeader(req.Header),
			"body":        *reqBody,
			"request_uri": req.URL.RequestURI(),
		})

		var level = 4
		var response []byte
		if err != nil {
			level = 1
			response, _ = json.Marshal(map[string]interface{}{
				"err": err.Error(),
			})
		} else {
			rawRespBody, err = io.ReadAll(resp.Body)
			resp.Body = io.NopCloser(bytes.NewBuffer(rawRespBody))

			respBody := new(interface{})
			if err := json.Unmarshal(rawRespBody, respBody); err != nil {
				*respBody = string(rawRespBody)
			}

			response, _ = json.Marshal(map[string]interface{}{
				"status_code": resp.StatusCode,
				"header":      convertHTTPHeader(resp.Header),
				"body":        *respBody,
			})
		}

		data := map[string]interface{}{
			"topic":        logTopic,
			"level":        level,
			"v1":           getV1(c.Context),
			"v2":           getV2(c.Context),
			"v3":           getV3(c.Context),
			"message":      "api日志:" + req.URL.String(),
			"request":      string(request),
			"response":     string(response),
			"create_time":  time.Now().Format("2006-01-02 15:04:05.000000"),
			"request_id":   getRequestID(c.Context),
			"execute_time": time.Since(start).Milliseconds(),
		}

		v, _ := json.Marshal(data)
		logger(v)
	}()

	if err != nil {
		return nil, err
	}

	return resp, nil
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

func getRequestID(ctx Context) string {
	if ctx == nil {
		return ""
	}
	value, exist := ctx.Get(requestID)
	if !exist {
		return ""
	}
	v, ok := value.(string)
	if !ok {
		return ""
	}
	return v
}

func getV1(ctx Context) string {
	if ctx == nil {
		return ""
	}
	value, exist := ctx.Get(v1)
	if !exist {
		return ""
	}
	v, ok := value.(string)
	if !ok {
		return ""
	}
	return v
}

func getV2(ctx Context) string {
	if ctx == nil {
		return ""
	}
	value, exist := ctx.Get(v2)
	if !exist {
		return ""
	}
	v, ok := value.(string)
	if !ok {
		return ""
	}
	return v
}

func getV3(ctx Context) string {
	if ctx == nil {
		return ""
	}
	value, exist := ctx.Get(v3)
	if !exist {
		return ""
	}
	v, ok := value.(string)
	if !ok {
		return ""
	}
	return v
}
