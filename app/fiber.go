package app

import (
	"flapypan-blog-go/tool"
	"github.com/gofiber/fiber/v2"
	"log"
)

// CreateApp 创建 fiber
func CreateApp() *fiber.App {
	app := fiber.New(fiber.Config{
		// 应用名称
		AppName: "FlapyPan Blog API",
		// 全局错误处理
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			errInfo := err.Error()
			log.Println(errInfo)
			return c.JSON(tool.Err(500, errInfo))
		},
	})
	return app
}
