package middleware

import (
	"flapypan-blog-go/tool"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Protected 认证中间件
func Protected() func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		// 获取 session
		sess := ctx.Locals("session").(session.Session)
		// 检查 session 中的登录状态
		if sess.Get("loginStatus") == true {
			// 如果登录则放行
			return ctx.Next()
		}
		return ctx.JSON(tool.ErrCode(401))
	}
}
