package data

import (
	"time"
)

type Illust struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	AuthorID    string       `json:"authorid"`
	Author      string       `json:"author"`
	Account     string       `json:"account"`
	R18         bool         `json:"r18"`
	CreateDate  time.Time    `json:"create_date"`
	UploadDate  time.Time    `json:"upload_date"`
	Tags        []string     `json:"tags"`
	Pages       []IllustPage `json:"pages"`
}

type IllustPage struct {
	P        int    `json:"id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Thumb    string `json:"thumb"`
	Small    string `json:"small"`
	Regular  string `json:"regular"`
	Original string `json:"original"`
}
