package db

import (
	"context"
	. "flapypan-blog-go/conf"
	"flapypan-blog-go/pagination"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

var DB *pgxpool.Pool

// ConnectDB 连接到数据库
func ConnectDB() {

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		Env("DB_USERNAME"), Env("DB_PASSWORD"),
		Env("DB_HOST"), Env("DB_PORT"),
		Env("DB_NAME"))
	log.Println("😚 尝试连接到数据库...")
	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("😅 数据库连接失败 %v", err)
	}
	// 检查连接是否成功
	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("😅 检查数据库连接失败 %v", err)
	}
	DB = db
	log.Println("😁 成功连接到数据库！")
}

type FieldMapper[T any] func(ele *T) []any

func MapToSlice[T any](rows *pgx.Rows, slice []*T, mapper FieldMapper[T]) ([]*T, error) {
	for (*rows).Next() {
		var ele T
		fields := mapper(&ele)
		if err := (*rows).Scan(fields...); err != nil {
			return nil, err
		}
		slice = append(slice, &ele)
	}
	return slice, nil
}

func MapToPage[T any](rows *pgx.Rows, pageable *pagination.Pageable, mapper FieldMapper[T]) (*pagination.Page[T], error) {
	slice := make([]*T, 0, pageable.Size)
	slice, err := MapToSlice(rows, slice, mapper)
	if err != nil {
		return nil, err
	}
	return pagination.NewPage(slice, pageable), nil
}
