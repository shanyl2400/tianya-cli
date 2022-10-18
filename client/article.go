package client

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	contentPageSize = 100
)

var (
	ErrPageOutRange = errors.New("page out of range")
	ErrNoHistory    = errors.New("no history")
)

type Article struct {
	Title      string
	Author     string
	ViewCount  int
	ReplyCount int

	ReplyAt time.Time

	Href string
	Type string

	posts []*Post
	index int

	content        string
	curViewContent string
	offset         int

	historyNext []string
	historyPrev []string

	prevPage string
	nextPage string
}

func (c *Article) Open() error {
	return c.open(c.Href)
}

func (c *Article) Next() (string, error) {
	if len(c.historyNext) > 0 {
		out := c.historyNext[len(c.historyNext)-1]
		c.historyNext = c.historyNext[:len(c.historyNext)-1]

		if c.curViewContent != "" {
			c.historyPrev = append(c.historyPrev, c.curViewContent)
		}
		c.curViewContent = out
		return out, nil
	}

	for c.content == "" {
		post, err := c.NextPost()
		if err != nil {
			return "", err
		}
		c.content = post.Content
		c.offset = 0
	}

	content := []rune(c.content)

	right := c.offset + contentPageSize
	if right >= len(content) {
		right = len(content) - 1
	}

	out := content[c.offset:right]
	c.offset = c.offset + contentPageSize

	if c.offset >= len(content) {
		c.content = ""
	}
	if c.curViewContent != "" {
		c.historyPrev = append(c.historyPrev, c.curViewContent)
	}
	c.curViewContent = string(out)
	return c.curViewContent, nil
}

func (c *Article) Prev() (string, error) {
	if len(c.historyPrev) > 0 {
		out := c.historyPrev[len(c.historyPrev)-1]
		c.historyPrev = c.historyPrev[:len(c.historyPrev)-1]

		if c.curViewContent != "" {
			c.historyNext = append(c.historyNext, c.curViewContent)
		}
		c.curViewContent = out
		return out, nil
	}

	return "", ErrNoHistory
}

func (c *Article) NextPost() (*Post, error) {
	if c.index >= len(c.posts) {
		err := c.NextPage()
		if err != nil {
			return nil, err
		}
		c.index = 0
		return c.NextPost()
	}
	post := c.posts[c.index]
	c.index++
	return post, nil
}

func (c *Article) NextPage() error {
	if c.nextPage == "" {
		return ErrPageOutRange
	}
	return c.open(c.nextPage)
}

func (c *Article) PrevPage() error {
	if c.prevPage == "" {
		return ErrPageOutRange
	}
	return c.open(c.prevPage)
}

func (c *Article) open(path string) error {
	res, err := http.Get(baseURL + path)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("access website failed, status: %v", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	// Find the review items
	doc.Find(".atl-con-bd .bbs-content").Each(func(i int, s *goquery.Selection) {
		content := s.Text()
		c.posts = append(c.posts, &Post{
			Content: strings.TrimSpace(content),
		})
	})

	c.prevPage = ""
	c.nextPage = ""
	doc.Find(".atl-pages form a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if s.Text() == "上页" {
			c.prevPage = href
		} else if s.Text() == "下页" {
			c.nextPage = href
		}
	})
	return nil
}

func NewArticle() *Article {
	return &Article{
		posts:       make([]*Post, 0),
		historyNext: make([]string, 0),
		historyPrev: make([]string, 0),
	}
}
