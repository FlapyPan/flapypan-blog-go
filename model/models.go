package model

import (
	"database/sql"
	"time"
)

type Access struct {
	ID         int64     `json:"id" toml:"id" yaml:"id"`
	IP         string    `json:"ip" toml:"ip" yaml:"ip"`
	Referrer   string    `json:"referrer" toml:"referrer" yaml:"referrer"`
	Ua         string    `json:"ua" toml:"ua" yaml:"ua"`
	CreateDate time.Time `json:"createDate" toml:"createDate" yaml:"createDate"`
	ArticleID  int64     `json:"articleId" toml:"articleId" yaml:"articleId"`
}

type Link struct {
	ID   int64  `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name string `boil:"name" json:"name" toml:"name" yaml:"name"`
	URL  string `boil:"url" json:"url" toml:"url" yaml:"url"`
}

type Setting struct {
	SKey   string         `json:"sKey" toml:"sKey" yaml:"sKey"`
	SValue sql.NullString `json:"sValue" toml:"sValue" yaml:"sValue"`
}
