package handler

import (
	"flapypan-blog-go/model"
	. "flapypan-blog-go/pagination"
	. "flapypan-blog-go/tool"
	"github.com/gofiber/fiber/v2"
	"time"
)

func ArticlePage(ctx *fiber.Ctx) error {
	pageable := NewPageable(ctx)
	dbCtx := ctx.UserContext()
	articlePage, err := model.QueryArticlePage(&dbCtx, pageable)
	if err != nil {
		return err
	}
	for _, a := range articlePage.Content {
		tags, err := model.QueryTagListByArticle(&dbCtx, a)
		if err != nil {
			return err
		}
		a.Tags = tags
	}
	return ctx.JSON(OkData(articlePage))
}

func ArticleYearCount(ctx *fiber.Ctx) error {
	dbCtx := ctx.UserContext()
	list, err := model.QueryArticleYearlyCountList(&dbCtx)
	if err != nil {
		return err
	}
	return ctx.JSON(OkData(list))
}

func ArticleListByYear(ctx *fiber.Ctx) error {
	dbCtx := ctx.UserContext()
	year, err := ctx.ParamsInt("year", time.Now().Year())
	if err != nil {
		return ctx.JSON(ErrCode(400))
	}
	articles, err := model.QueryArticleListByYear(&dbCtx, year)
	if err != nil {
		return err
	}
	for _, a := range articles {
		tags, err := model.QueryTagListByArticle(&dbCtx, a)
		if err != nil {
			return err
		}
		a.Tags = tags
	}
	return ctx.JSON(OkData(articles))
}

func ArticleByPath(ctx *fiber.Ctx) error {
	path := ctx.Params("path", "")
	if !ValidArticlePath(path) {
		return ctx.JSON(ErrCode(400))
	}
	dbCtx := ctx.UserContext()
	article, err := model.QueryArticleByPath(&dbCtx, path)
	if err != nil {
		return err
	}
	if article.ID <= 0 {
		return ctx.JSON(ErrCode(404))
	}
	tags, err := model.QueryTagListByArticle(&dbCtx, article)
	if err != nil {
		return err
	}
	article.Tags = tags
	return ctx.JSON(OkData(article))
}

func AddArticle(ctx *fiber.Ctx) error {
	var saveReq model.ArticleSaveReq
	err := ctx.BodyParser(&saveReq)
	if err != nil {
		return err
	}
	err = GetValidateError(&saveReq)
	if err != nil {
		return err
	}
	dbCtx := ctx.UserContext()
	path, err := model.AddArticle(&dbCtx, &saveReq)
	if err != nil {
		return err
	}
	return ctx.JSON(OkData(path))
}

func ModifyArticle(ctx *fiber.Ctx) error {
	var saveReq model.ArticleSaveReq
	err := ctx.BodyParser(&saveReq)
	if err != nil {
		return err
	}
	if saveReq.ID <= 0 {
		return ctx.JSON(ErrCode(400))
	}
	err = GetValidateError(&saveReq)
	if err != nil {
		return err
	}
	dbCtx := ctx.UserContext()
	path, err := model.ModifyArticle(&dbCtx, &saveReq)
	if err != nil {
		return err
	}
	return ctx.JSON(OkData(path))
}

func DeleteArticle(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id", time.Now().Year())
	if err != nil {
		return ctx.JSON(ErrCode(400))
	}
	dbCtx := ctx.UserContext()
	err = model.LogicDeleteArticle(&dbCtx, int64(id))
	if err != nil {
		return err
	}
	return ctx.JSON(Ok())
}
