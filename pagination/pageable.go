package pagination

import "github.com/gofiber/fiber/v2"

type Pageable struct {
	Keyword    string `query:"keyword" json:"keyword,omitempty"`
	Page       int    `query:"page" json:"page,omitempty"`
	Size       int    `query:"size" json:"size,omitempty"`
	Offset     int    `query:"-" json:"offset,omitempty"`
	Total      int64  `query:"-" json:"total,omitempty"`
	TotalPages int64  `query:"-" json:"totalPages,omitempty"`
}

// NewPageable 获取分页信息
func NewPageable(ctx *fiber.Ctx) *Pageable {
	// 获取分页参数
	pageable := &Pageable{}
	if ctx.QueryParser(pageable) != nil {
		return nil
	}
	// 修正页数
	if pageable.Page < 0 {
		pageable.Page = 0
	}
	// 限制每页数量
	switch {
	case pageable.Size > 100:
		pageable.Size = 100
	case pageable.Size <= 0:
		pageable.Size = 6
	}
	// 计算偏移量
	pageable.Offset = pageable.Page * pageable.Size
	return pageable
}

// SetTotal 计算总页数
func (p *Pageable) SetTotal(total int64) *Pageable {
	p.Total = total
	totalPages := p.Total / int64(p.Size)
	if p.Total%int64(p.Size) != 0 {
		totalPages += 1
	}
	p.TotalPages = totalPages
	return p
}

type Page[T any] struct {
	Size       int   `json:"size,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int64 `json:"totalPages,omitempty"`
	Content    []*T  `json:"content,omitempty"`
}

func NewPage[T any](content []*T, pageable *Pageable) *Page[T] {
	return &Page[T]{
		Size:       pageable.Size,
		Total:      pageable.Total,
		TotalPages: pageable.TotalPages,
		Content:    content,
	}
}
