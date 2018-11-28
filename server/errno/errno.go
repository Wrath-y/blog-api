package errno

import (
	"bytes"
)

type Errno struct {
	Code int
	Message string
}

func (err Errno) Error() string {
	return err.Message
}

type Err struct {
	Code int
	Message string
	Err string
}

func New(errno *Errno, err string) *Err {
	return &Err{Code: errno.Code, Message: errno.Message, Err:err}
}

func (err *Err) Add(message string) error {
	err.Message += " " + message

	return err
}

func (err *Err) AddParam(param string) error {
	var buffer bytes.Buffer
	buffer.WriteString(err.Message)
	buffer.WriteString(" ")
	buffer.WriteString(param)
	err.Message = buffer.String()

	return err
}

func (err *Err) Error() string {
	var buffer bytes.Buffer
	buffer.WriteString(err.Message)
	buffer.WriteString(err.Err)

	return buffer.String()
}

func ReturnErr(err error) (int, string) {
	if err == nil {
		return Success.Code, Success.Message
	}

	switch typed := err.(type) {
	case *Err:
		return typed.Code, typed.Error()
	case *Errno:
		return typed.Code, typed.Message
	default:
	}

	return ServerError.Code, err.Error()
}
