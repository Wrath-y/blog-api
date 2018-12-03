package errno

var (
	Success	= &Errno{Code: 0, Message: "Success"}
	ServerError = &Errno{Code: 10001, Message: "Server Error"}
	RouteError = &Errno{Code: 10002, Message: "Route Error"}
	DatabaseError = &Errno{Code: 10003, Message: "Database Error"}
	UploadError = &Errno{Code: 10004, Message: "Upload Error"}
	BindError = &Errno{Code: 20001, Message: "Bind Error"}
	RequestError = &Errno{Code: 20002, Message: "Request Error"}
)