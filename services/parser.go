package services

import (
	"time"

	"github.com/aitsuki/pixiv-captor/data"
)

type IllustData struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	UserId      string     `json:"userId"`
	UserName    string     `json:"userName"`
	UserAccount string     `json:"userAccount"`
	CreateDate  time.Time  `json:"createDate"`
	UploadDate  time.Time  `json:"uploadDate"`
	Tags        TagsData   `json:"tags"`
	Pages       []PageData `json:"pages"`
}

func (a *IllustData) isR18() bool {
	for _, t := range a.Tags.Tags {
		if t.Tag == "R-18" {
			return true
		}
	}
	return false
}

func (a *IllustData) stringTags() []string {
	tags := make([]string, 0, len(a.Tags.Tags))
	for _, tag := range a.Tags.Tags {
		tags = append(tags, tag.Tag)
		if tag.Translation != nil {
			tags = append(tags, tag.Translation.En)
		}
	}
	return tags
}

func (a *IllustData) pages() []data.IllustPage {
	pages := make([]data.IllustPage, 0, len(a.Pages))
	for i, p := range a.Pages {
		pages = append(pages, p.toPageEntity(a.ID, i))
	}
	return pages
}

type TagsData struct {
	Tags []TagData `json:"tags"`
}

type TagData struct {
	Tag         string           `json:"tag"`
	Translation *TranslationData `json:"translation,omitempty"`
}

type TranslationData struct {
	En string `json:"en"`
}

type PageData struct {
	Urls   UrlsData `json:"urls"`
	Width  int      `json:"width"`
	Height int      `json:"height"`
}

func (p *PageData) toPageEntity(illustID string, i int) data.IllustPage {
	return data.IllustPage{
		P:        i,
		Width:    p.Width,
		Height:   p.Height,
		Thumb:    p.Urls.ThumbMini,
		Small:    p.Urls.Small,
		Regular:  p.Urls.Regular,
		Original: p.Urls.Original,
	}
}

type UrlsData struct {
	ThumbMini string `json:"thumb_mini"`
	Small     string `json:"small"`
	Regular   string `json:"regular"`
	Original  string `json:"original"`
}

func (a *IllustData) ToEntity() *data.Illust {
	return &data.Illust{
		ID:          a.ID,
		Title:       a.Title,
		Description: a.Description,
		AuthorID:    a.UserId,
		Author:      a.UserName,
		Account:     a.UserAccount,
		R18:         a.isR18(),
		CreateDate:  a.CreateDate,
		UploadDate:  a.UploadDate,
		Tags:        a.stringTags(),
		Pages:       a.pages(),
	}
}
