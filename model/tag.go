package model

import (
	"context"
	. "flapypan-blog-go/db"
)

type Tag struct {
	ID   int64  `json:"id,omitempty" toml:"id" yaml:"id"`
	Name string `json:"name,omitempty" toml:"name" yaml:"name"`
}

func tagMapper(tag *Tag) []any {
	return []any{
		&tag.ID,
		&tag.Name,
	}
}

const sqlQueryTagListByArticle = `
select tag_id as id, name
from t_tag join t_article_tag on t_tag.id = t_article_tag.tag_id
where article_id = $1`

// QueryTagListByArticle 通过文章查询标签
func QueryTagListByArticle(ctx *context.Context, article *Article) ([]*Tag, error) {
	rows, err := DB.Query(*ctx, sqlQueryTagListByArticle, article.ID)
	if err != nil {
		return nil, err
	}
	tags := make([]*Tag, 0)
	return MapToSlice(&rows, tags, tagMapper)
}
