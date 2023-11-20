package router

import (
	"flapypan-blog-go/conf"
	. "flapypan-blog-go/handler"
	. "flapypan-blog-go/middleware"
	"github.com/gofiber/fiber/v2"
	"time"
)

// 静态文件处理
func static(app *fiber.App) {
	app.Static("/static", conf.Env("UPLOAD_DIR"), fiber.Static{
		// 开启压缩
		Compress: true,
		// 非下载模式
		Download: false,
		// 服务器缓存 3 小时
		CacheDuration: 3 * time.Hour,
		// 客户端缓存 1 年
		MaxAge: 31536000,
	})
}

// 认证相关路由
func authRouter(app *fiber.App) {
	app.Group("/auth").
		Use(Session()).
		Get("/", Protected(), CheckLogin).
		Post("/", Login).
		Get("/logout", Logout)
}

// 文章相关路由
func articleRouter(app *fiber.App) {
	app.Group("/article").
		Get("/", ArticlePage).
		Get("/year", ArticleYearCount).
		Get("/year/:year", ArticleListByYear).
		Get("/:path", ArticleByPath)
	app.Group("/article").
		Use(Session(), Protected()).
		Post("/", AddArticle).
		Put("/", ModifyArticle).
		Delete("/:id", DeleteArticle)
}

// Setup 设置所有的路由
func Setup(app *fiber.App) {
	static(app)
	authRouter(app)
	articleRouter(app)
}
