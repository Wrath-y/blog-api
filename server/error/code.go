package error

var (
	Success	= &Errno{Code: 0, Message: "Success"}
	ServerError = &Errno{Code: 10001, Message: "Server Error"}
)