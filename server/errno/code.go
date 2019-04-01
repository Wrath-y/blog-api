package errno

var (
	Success	= &Errno{Code: 0, Message: "Success"}
	ServerError = &Errno{Code: 10001, Message: "Server Error"}
	RouteError = &Errno{Code: 10002, Message: "Route Error"}
	DatabaseError = &Errno{Code: 10003, Message: "Database Error"}
	UploadError = &Errno{Code: 10004, Message: "Upload Error"}
	BindError = &Errno{Code: 20001, Message: "Bind Error"}
	RequestError = &Errno{Code: 20002, Message: "Request Error"}
	ErrToken = &Errno{Code: 20010, Message: "Token Error"}
	ErrUserNotFound = &Errno{Code: 2011, Message: "User not found"}
	ErrPasswordIncorrect = &Errno{Code: 2012, Message: "Password incorrent"}
	ErrTokenInvalid = &Errno{Code: 2013, Message: "Token invalid"}
	ErrCurl = &Errno{Code: 2014, Message: "Curl Error"}
	ErrIoCopy = &Errno{Code: 2015, Message: "io.Copy Error"}
	ErrOsCreate = &Errno{Code: 2016, Message: "os.Create Error"}
	ErrIoutilReadAll = &Errno{Code: 2017, Message: "ioutil.ReadAll Error"}
	ErrExp = &Errno{Code: 2018, Message: "Exp Error"}
)