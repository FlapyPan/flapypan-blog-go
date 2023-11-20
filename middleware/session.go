package middleware

import (
	"flapypan-blog-go/tool"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"time"
)

// SessionStore 管理会话存储
var sessionStore = session.New(session.Config{
	Expiration:     24 * time.Hour,
	CookieSecure:   true,
	CookieSameSite: "same-site",
})

// Session session 中间件
func Session() func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		// 获取 session
		sess, err := sessionStore.Get(ctx)
		if err != nil {
			return ctx.JSON(tool.ErrCode(500))
		}
		ctx.Locals("session", sess)
		return ctx.Next()
	}
}
