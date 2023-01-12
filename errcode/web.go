package errcode

// 微信端错误码
var (
	WebNetworkBusy  = &ErrCode{40000, "网络繁忙，请稍后重试", ""}
	WebInvalidParam = &ErrCode{40100, "无效的参数", ""}
	WebInvalidSign  = &ErrCode{40101, "无效的签名", ""}
	WebBodyTooLarge = &ErrCode{40102, "请求消息体过大", ""}
)
