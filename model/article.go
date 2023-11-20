package model

import (
	"context"
	. "flapypan-blog-go/db"
	"flapypan-blog-go/pagination"
	"time"
)

type Article struct {
	ID         int64     `json:"id,omitempty" toml:"id" yaml:"id"`
	Title      string    `json:"title,omitempty" toml:"title" yaml:"title"`
	Path       string    `json:"path,omitempty" toml:"path" yaml:"path"`
	Cover      *string   `json:"cover,omitempty" toml:"cover" yaml:"cover"`
	Content    string    `json:"content,omitempty" toml:"content" yaml:"content"`
	CreateDate time.Time `json:"createDate,omitempty" toml:"createDate" yaml:"createDate"`
	UpdateDate time.Time `json:"updateDate,omitempty" toml:"updateDate" yaml:"updateDate"`
	Deleted    bool      `json:"deleted,omitempty" toml:"deleted" yaml:"deleted"`

	Tags []*Tag `json:"tags" toml:"tags" yaml:"tags"`
}

func articleMapper(article *Article) []any {
	return []any{
		&article.ID,
		&article.Title,
		&article.Path,
		&article.Cover,
		&article.Content,
		&article.CreateDate,
		&article.UpdateDate,
	}
}

func articleMapperWithoutContent(article *Article) []any {
	return []any{
		&article.ID,
		&article.Title,
		&article.Path,
		&article.Cover,
		&article.CreateDate,
		&article.UpdateDate,
	}
}

// QueryArticlePage 查询文章分页
func QueryArticlePage(ctx *context.Context, pageable *pagination.Pageable) (*pagination.Page[Article], error) {
	if pageable.Keyword == "" {
		return articlePage(ctx, pageable)
	} else {
		return articlePageByKeyword(ctx, pageable)
	}
}

const sqlArticleTotalPages = `select count(*) as total from t_article where deleted = false`
const sqlArticlePage = `
select distinct t_article.id as id, title, path, cover, create_date, update_date
from t_article
where deleted = false
order by update_date desc, create_date desc
limit $1 offset $2`

func articlePage(ctx *context.Context, pageable *pagination.Pageable) (*pagination.Page[Article], error) {
	var total int64
	if err := DB.QueryRow(*ctx, sqlArticleTotalPages).Scan(&total); err != nil {
		return nil, err
	}
	pageable.SetTotal(total)
	rows, err := DB.Query(*ctx, sqlArticlePage, pageable.Size, pageable.Offset)
	if err != nil {
		return nil, err
	}
	return MapToPage(&rows, pageable, articleMapperWithoutContent)
}

const sqlArticleTotalPagesByKeyword = `select count(distinct t_article.id)
from t_article
         join t_article_tag on t_article_tag.article_id = t_article.id
         join t_tag on t_tag.id = t_article_tag.tag_id
where deleted = false
  and (lower(title) like $1 or t_tag.name like $2)`
const sqlArticlePageByKeyword = `
select distinct t_article.id as id, title, path, cover, create_date, update_date
from t_article
         join t_article_tag on t_article_tag.article_id = t_article.id
         join t_tag on t_tag.id = t_article_tag.tag_id
where deleted = false
  and (lower(title) like $1 or t_tag.name like $2)
order by update_date desc, create_date desc
limit $3 offset $4`

func articlePageByKeyword(ctx *context.Context, pageable *pagination.Pageable) (*pagination.Page[Article], error) {
	keyword := "%" + pageable.Keyword + "%"
	var total int64
	err := DB.QueryRow(*ctx, sqlArticleTotalPagesByKeyword, keyword, keyword).Scan(&total)
	if err != nil {
		return nil, err
	}
	pageable.SetTotal(total)
	rows, err := DB.Query(*ctx, sqlArticlePageByKeyword, keyword, keyword, pageable.Size, pageable.Offset)
	if err != nil {
		return nil, err
	}
	return MapToPage(&rows, pageable, articleMapperWithoutContent)
}

const sqlArticleTotalPagesByTag = `
select count(distinct t_article.id)
from t_article
         join t_article_tag on t_article_tag.article_id = t_article.id
         join t_tag on t_tag.id = t_article_tag.tag_id
where deleted = false
  and t_tag.name = $1`
const sqlArticlePageByTag = `
select distinct t_article.id as id, title, path, cover, create_date, update_date
from t_article
         join t_article_tag on t_article_tag.article_id = t_article.id
         join t_tag on t_tag.id = t_article_tag.tag_id
where deleted = false
  and t_tag.name = $1
order by update_date desc, create_date desc
limit $2 offset $3`

// QueryArticlePageByTag 通过标签查询文章分页
func QueryArticlePageByTag(ctx *context.Context, pageable *pagination.Pageable, tagName string) (*pagination.Page[Article], error) {
	var total int64
	err := DB.QueryRow(*ctx, sqlArticleTotalPagesByTag).Scan(&total)
	if total == 0 && err != nil {
		return nil, err
	}
	pageable.SetTotal(total)
	rows, err := DB.Query(*ctx, sqlArticlePageByTag, tagName, pageable.Size, pageable.Offset)
	if err != nil {
		return nil, err
	}
	return MapToPage(&rows, pageable, articleMapperWithoutContent)
}

type ArticleYearlyCount struct {
	Year  int   `json:"year,omitempty"`
	Count int64 `json:"count,omitempty"`
}

func articleYearlyCountMapper(ele *ArticleYearlyCount) []any {
	return []any{&ele.Year, &ele.Count}
}

const sqlArticleYearlyCount = `
select extract(year from create_date) as year, count(id) as count
from t_article
where deleted = false
group by year
order by year desc`

// QueryArticleYearlyCountList 获取每个年份下文章的数量
func QueryArticleYearlyCountList(ctx *context.Context) ([]*ArticleYearlyCount, error) {
	rows, err := DB.Query(*ctx, sqlArticleYearlyCount)
	if err != nil {
		return nil, err
	}
	var yearlyCounts []*ArticleYearlyCount
	return MapToSlice(&rows, yearlyCounts, articleYearlyCountMapper)
}

const sqlArticleListByYear = `
select id, path, title, create_date
from t_article
where deleted = false
  and extract(year from create_date) = $1`

// QueryArticleListByYear 查询指定年份下的文章
func QueryArticleListByYear(ctx *context.Context, year int) ([]*Article, error) {
	rows, err := DB.Query(*ctx, sqlArticleListByYear, year)
	if err != nil {
		return nil, err
	}
	slice := make([]*Article, 0)
	return MapToSlice(&rows, slice, func(ele *Article) []any {
		return []any{&ele.ID, &ele.Path, &ele.Title, &ele.CreateDate}
	})
}

const sqlQueryArticleByPath = `
select id, title, path, cover, content, create_date, update_date
from t_article
where deleted = false
  and path = $1`

func QueryArticleByPath(ctx *context.Context, path string) (*Article, error) {
	var article Article
	row := DB.QueryRow(*ctx, sqlQueryArticleByPath, path)
	if err := row.Scan(articleMapper(&article)...); err != nil {
		return nil, err
	}
	return &article, nil
}

type ArticleSaveReq struct {
	ID       int64    ` json:"id,omitempty"`
	Title    string   `json:"title,omitempty" validate:"required,min=2,max=32"`
	Path     string   `json:"path,omitempty"  validate:"required,valid-path"`
	Cover    *string  `json:"cover,omitempty"`
	Content  string   `json:"content,omitempty" validate:"required"`
	TagNames []string `json:"tagNames,omitempty" toml:"tagNames" yaml:"tagNames"`
}

const sqlAddArticle = `
insert into t_article(title, path, cover, content, create_date, update_date)
values ($1, $2, $3, $4, current_timestamp, current_timestamp)
returning path`

// AddArticle 添加文章，返回访问路径
func AddArticle(ctx *context.Context, article *ArticleSaveReq) (string, error) {
	exec, err := DB.Exec(*ctx, sqlAddArticle, article.Title, article.Path, article.Cover, article.Content)
	if err != nil {
		return "", err
	}
	// todo 保存标签
	return exec.String(), nil
}

const sqlModifyArticle = `
update t_article
set title       = $1,
    path        = $2,
    cover       = $3,
    content     = $4,
    update_date = current_timestamp
where id = $5
returning path`

// ModifyArticle 修改文章，返回访问路径
func ModifyArticle(ctx *context.Context, article *ArticleSaveReq) (string, error) {
	exec, err := DB.Exec(*ctx, sqlModifyArticle, article.Title, article.Path, article.Cover, article.Content, article.ID)
	if err != nil {
		return "", err
	}
	// todo 保存标签
	return exec.String(), nil
}

const sqlLogicDeleteArticle = `
update t_article
set deleted     = true,
    update_date = current_timestamp
where id = $1`

// LogicDeleteArticle 逻辑删除文章
func LogicDeleteArticle(ctx *context.Context, id int64) error {
	_, err := DB.Exec(*ctx, sqlLogicDeleteArticle, id)
	return err
}
