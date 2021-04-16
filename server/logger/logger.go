package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type fileWriter struct {
	logRoot  string
	fileName string

	f                   *os.File
	mutex               sync.RWMutex
	lastUpdateHourInDay int64
}


func NewFileWriter(path, name string) *fileWriter {
	return &fileWriter{logRoot: path, fileName: name}
}

func (self *fileWriter) Write(p []byte) (int, error) {
	now := time.Now()
	hoursInDay := now.Unix() / 3600
	self.mutex.Lock()
	defer self.mutex.Unlock()
	if self.f == nil || hoursInDay != self.lastUpdateHourInDay { // 15:04:05
		filename := fmt.Sprintf("%s%s.log", self.fileName, now.Format("20060102"))
		f, err := os.OpenFile(fmt.Sprintf("%s%s", self.logRoot, filename), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600|0644)

		if err != nil {
			return 0, err
		}
		self.f = f
		self.lastUpdateHourInDay = hoursInDay
	}
	n, err := self.f.Write(p)
	if err != nil {
		self.f = nil
	}
	return n, err
}