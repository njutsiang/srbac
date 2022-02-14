package code

type Code int

// 错误码
var (
	ParamsError Code = 400000
	CsrfTokenError Code = 400001
	UserNotLogin Code = 401000
	UserNotExists Code = 404001
	UnknownModelFieldType Code = 500100
)

// 错误提示语
var messages = map[Code]string{
	ParamsError: "参数错误",
	CsrfTokenError: "表单验证失败",
	UserNotExists: "用户不存在",
	UserNotLogin: "用户未登录",
	UnknownModelFieldType: "未知的模型字段数据类型",
}

// 获取错误提示语
func GetMessage(code Code) string {
	return messages[code]
}
