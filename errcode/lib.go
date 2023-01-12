package errcode

// 框架组件错误码
var (
	LibNoRoute   = &ErrCode{1001, "路由未找到", ""}
	LibRateLimit = &ErrCode{1002, "服务繁忙，请稍后再试", ""}
)
