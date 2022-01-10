// Package code /**
package code

type ResponseCode int

const (
	SuccessCode ResponseCode = 0

	InvalidAuth ResponseCode = 401

	BadRequest     ResponseCode = 1000
	InvalidRequest ResponseCode = 1000 + iota
	InvalidParams

	InvalidPermission = 4001
	AccountDisabled   = 4006
)

var ResponseCodeMap = map[ResponseCode]string{
	SuccessCode:       "请求成功",
	InvalidAuth:       "身份校验失败",
	InvalidRequest:    "无效的请求",
	BadRequest:        "系统错误",
	InvalidParams:     "请求参数无效",
	InvalidPermission: "权限不足",
	AccountDisabled:   "账号未激活或已禁用",
}

func GetCodeMsg(code ResponseCode) string {
	h, ok := ResponseCodeMap[code]
	if !ok {
		return "未知错误"
	}
	return h
}
