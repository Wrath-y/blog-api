package errno

var (
	Success	= &Errno{Code: 0, Message: "Success"}
	ServerError = &Errno{Code: 10001, Message: "Server Error"}
	BindError = &Errno{Code: 10002, Message: "Bind Error"}
	TitleError = &Errno{Code: 20001, Message: "Title Error"}
)