package errno

var (
	Success              = &Errno{Code: 0, Message: "success"}
	ServerError          = &Errno{Code: 10001, Message: "server error"}
	RouteError           = &Errno{Code: 10002, Message: "route error"}
	DatabaseError        = &Errno{Code: 10003, Message: "database error"}
	BindError            = &Errno{Code: 20000, Message: "bind error"}
	RequestError         = &Errno{Code: 20001, Message: "request error"}
	TokenErr             = &Errno{Code: 20002, Message: "token error"}
	UserNotFoundErr      = &Errno{Code: 30000, Message: "user not found"}
	PasswordIncorrectErr = &Errno{Code: 30001, Message: "password incorrent"}
	TokenInvalidErr      = &Errno{Code: 30002, Message: "token invalid"}
	CurlErr              = &Errno{Code: 30003, Message: "curl error"}
	RegexpErr            = &Errno{Code: 30004, Message: "regexp error"}
	IndexOutOfRangeErr   = &Errno{Code: 30005, Message: "index out of range"}
	UploadErr            = &Errno{Code: 40000, Message: "upload error"}
)
