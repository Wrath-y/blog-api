package error

import "fmt"

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
	Err error
}

func New(errno *Errno, err error) *Err {
	return &Err{Code: errno.Code, Message: errno.Message, Err:err}
}

func (err *Err) Add(message string) error {
	err.Message += " " + message

	return err
}

func (err *Err) AddParams(format string, args ...interface{}) error {
	err.Message += " " + fmt.Sprintf(format, args...)

	return err
}

func (err *Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.Err)
}

func printErr(err error) (int, string) {
	if err == nil {
		return Success.Code, Success.Message
	}

	switch typed := err.(type) {
	case *Err:
		return typed.Code, typed.Message
	case *Errno:
		return typed.Code, typed.Message
	default:
	}

	return ServerError.Code, err.Error()
}
