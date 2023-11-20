package handler

import (
	. "flapypan-blog-go/conf"
	. "flapypan-blog-go/tool"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// 检查 bcrypt 加密过的密码
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Login 登录
func Login(ctx *fiber.Ctx) error {
	var err error
	// 获取 session
	ss := ctx.Locals("session").(session.Session)
	loginRequest := LoginRequest{}
	// 获取登录请求体
	if err = ctx.BodyParser(&loginRequest); err != nil {
		return err
	}
	// 校验请求
	if err = GetValidateError(loginRequest); err != nil {
		return err
	}
	// 检查密码
	if loginRequest.Username != Env("ADMIN_USERNAME") ||
		!checkPasswordHash(loginRequest.Password, Env("ADMIN_PASSWORD")) {
		return fmt.Errorf("用户名或密码错误")
	}
	// 将登录状态存入 session
	ss.Set("loginStatus", true)
	if err = ss.Save(); err != nil {
		return err
	}
	return ctx.JSON(Ok())
}

// CheckLogin 检查登录状态
func CheckLogin(ctx *fiber.Ctx) error {
	// 获取 session
	ss := ctx.Locals("session").(session.Session)
	// 检查 session 中的登录状态
	status := ss.Get("loginStatus") == true
	return ctx.JSON(OkData(status))
}

// Logout 登出
func Logout(ctx *fiber.Ctx) error {
	// 获取 session
	ss := ctx.Locals("session").(session.Session)
	// 删除 session 中的登录状态
	ss.Delete("loginStatus")
	if err := ss.Save(); err != nil {
		return err
	}
	return ctx.JSON(Ok())
}
