package data

import (
	"time"
)

type Illust struct {
	ID          string `gorm:"primaryKey"`
	Title       string
	Description string
	AuthorID    string
	Author      string
	Account     string
	R18         bool
	CreateDate  time.Time
	UploadDate  time.Time
	Tags        []string
	Pages       []IllustPage
}

type IllustPage struct {
	IllustID string `gorm:"primaryKey"`
	P        int    `gorm:"primaryKey"`
	Width    int
	Height   int
	Thumb    string
	Small    string
	Regular  string
	Original string
}
