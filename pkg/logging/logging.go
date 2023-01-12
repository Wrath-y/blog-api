package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
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
	topic      string
	logfn      logFunc
	mu         sync.Mutex
	once       = &sync.Once{}
	lastHandle *os.File
	filePath   = "docs/log/wrath.cc/"
	current    = filePath + "access.log"
	renameFmt  = filePath + "access.20060102150405.log"
	logHandle  = &logrus.Logger{
		Formatter: &SimpleFormatter{},
		Level:     logrus.InfoLevel,
	}
)

type SimpleFormatter struct{}

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
	writeToFile()
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
	case nil:
		return ""
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

func (*SimpleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	b.WriteString(fmt.Sprintf("%s\n", entry.Message))
	return b.Bytes(), nil
}

func FileLogger(data []byte) {
	logHandle.Printf("%s", string(data))
}

func StdoutLogger(data []byte) {
	var logMsg LogMsg
	json.Unmarshal(data, &logMsg) // nolint
	if logMsg.Level == InfoLevel {
		fmt.Fprintf(os.Stdout, string(data)+"\n\n")
	} else {
		fmt.Fprintf(os.Stdout, color.RedString(string(data))+"\n\n")
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

func writeToFile() {
	if viper.GetString("app.log.output") != "file" {
		return
	}

	once.Do(func() {
		file, err := os.OpenFile(current, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("打开日志文件失败：", err)
		}
		lastHandle = file
		logHandle.SetOutput(file)

		go func() {
			tick := time.Tick(time.Second * 7)
			for t := range tick {
				info, _ := lastHandle.Stat()
				if info.Size() > 300<<20 {
					os.Rename(current, t.Format(renameFmt)) //nolint
					f, err := os.OpenFile(current, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
					if err == nil {
						logHandle.SetOutput(f)
						lastHandle.Close()
						lastHandle = f
					}
				}
				dirs, _ := os.ReadDir(filePath)
				if len(dirs) > 3 { // 超过3个删除最旧的
					os.Remove(filePath + dirs[0].Name())
				}
			}
		}()
	})
}
