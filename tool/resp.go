package tool

import "github.com/gofiber/fiber/v2"

// Resp 统一返回结构体
func Resp(code int, data interface{}, success bool) *fiber.Map {
	return &fiber.Map{
		"code":    code,
		"data":    data,
		"success": success,
	}
}

// Ok 成功响应
func Ok() *fiber.Map {
	return Resp(200, nil, true)
}

// OkData 成功响应，返回数据
func OkData(data interface{}) *fiber.Map {
	return Resp(200, data, true)
}

// Err 失败响应，返回状态码和数据
func Err(code int, data interface{}) *fiber.Map {
	return Resp(code, data, false)
}

// ErrCode 失败响应，返回状态码
func ErrCode(code int) *fiber.Map {
	return Resp(code, (interface{})(nil), false)
}

// ErrData 失败响应，返回数据
func ErrData(data interface{}) *fiber.Map {
	return Resp(400, data, false)
}
