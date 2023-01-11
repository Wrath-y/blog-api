package logging

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type level int8

const (
	FatalLevel level = iota + 1
	ErrorLevel
	WarnLevel
	InfoLevel
)

var (
	topic string
	logfn logFunc
	mu    sync.Mutex
)

type logFunc func([]byte)

type logger struct {
	v1        string
	v2        string
	v3        string
	requestID string
}

type LogMsg struct {
	Topic       string      `json:"topic"`
	Level       level       `json:"level"`
	RequestID   string      `json:"request_id"`
	V1          string      `json:"v1"`
	V2          string      `json:"v2"`
	V3          string      `json:"v3"`
	Message     interface{} `json:"message"`
	Request     interface{} `json:"request"`
	Response    interface{} `json:"response"`
	CreateTime  string      `json:"create_time"`
	ExecuteTime int64       `json:"execute_time"`
}

func Setup(t string, fn logFunc) {
	mu.Lock()
	defer mu.Unlock()
	if logfn != nil {
		return
	}

	topic = t
	logfn = fn
}

func New() *logger {
	return &logger{}
}

func (l *logger) SetRequestID(requestID string) {
	l.requestID = requestID
}

func (l *logger) SetV1(v1 string) {
	l.v1 = v1
}

func (l *logger) SetV2(v2 string) {
	l.v2 = v2
}

func (l *logger) SetV3(v3 string) {
	l.v3 = v3
}

func covert(val interface{}) interface{} {
	switch v := val.(type) {
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	case []byte:
		return string(v)
	default:
		return v
	}
}

func (l *logger) NewLogData(lv level, message string, request, response interface{}, t ...time.Time) []byte {
	file, line, ok := getFilterCallers()
	if ok {
		l.v2 = file + ":" + strconv.Itoa(line)
	}

	data := &LogMsg{
		Topic:      topic,
		Level:      lv,
		RequestID:  l.requestID,
		V1:         l.v1,
		V2:         l.v2,
		V3:         l.v3,
		Message:    message,
		Request:    covert(request),
		Response:   covert(response),
		CreateTime: time.Now().Format("2006-01-02 15:04:05.000000"),
	}
	if len(t) > 0 {
		data.ExecuteTime = time.Since(t[0]).Milliseconds()
	}

	b, _ := json.Marshal(data)
	return b
}

func (l *logger) Info(message string, request, response interface{}, t ...time.Time) {
	data := l.NewLogData(InfoLevel, message, request, response, t...)
	if logfn != nil {
		logfn(data)
	}
}

func (l *logger) Warn(message string, request, response interface{}, t ...time.Time) {
	data := l.NewLogData(WarnLevel, message, request, response, t...)
	if logfn != nil {
		logfn(data)
	}
}

func (l *logger) ErrorL(message string, request, response interface{}, t ...time.Time) {
	data := l.NewLogData(ErrorLevel, message, request, response, t...)
	if logfn != nil {
		logfn(data)
	}
}

func (l *logger) Fatal(message string, request, response interface{}, t ...time.Time) {
	data := l.NewLogData(FatalLevel, message, request, response, t...)
	if logfn != nil {
		logfn(data)
	}
}

// getFilterCallers 获取调用栈
func getFilterCallers() (file string, line int, ok bool) {
	for i := 2; i < 6; i++ {
		_, file, line, ok = runtime.Caller(i)
		if !ok {
			continue
		}

		return file, line, ok
	}
	return
}
