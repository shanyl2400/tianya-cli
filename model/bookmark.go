package model

import (
	"time"
)

type Bookmark struct {
	Title      string
	Author     string
	ViewCount  int
	ReplyCount int

	ReplyAt time.Time

	Href string
	Type string

	Posts []string
	Index int

	Content        string
	CurViewContent string
	Offset         int

	PrevPage string
	NextPage string
}
