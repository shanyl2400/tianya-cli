package client

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

type Column struct {
	id   string
	name string
	href string

	articleMap map[string]*Article
	nextPage   string
	prevPage   string
}

func (c *Column) Open() error {
	return c.enterPage(c.href)
}

func (c *Column) Name() string {
	return c.name
}

func (c *Column) NextPage() error {
	if c.nextPage == "" {
		return errors.New("no next page")
	}
	return c.enterPage(c.nextPage)
}

func (c *Column) PrevPage() error {
	if c.prevPage == "" {
		return errors.New("no prev page")
	}
	return c.enterPage(c.prevPage)
}

func (c *Column) enterPage(path string) error {
	res, err := http.Get(baseURL + path)
	if err != nil {
		log.Fatal(err)
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
	doc.Find(".tab-bbs-list tr").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			article := NewArticle()
			s.Children().Each(func(i int, s *goquery.Selection) {
				switch i {
				case 0:
					s.Find("span").Each(func(i int, s *goquery.Selection) {
						if i == 0 {
							title, _ := s.Attr("title")
							article.Type = title
						}
					})
					s.Find("a").Each(func(i int, s *goquery.Selection) {
						href, _ := s.Attr("href")
						article.Href = href
					})
					article.Title = strings.TrimSpace(s.Text())
				case 1:
					article.Author = strings.TrimSpace(s.Text())
				case 2:
					article.ViewCount, _ = strconv.Atoi(strings.TrimSpace(s.Text()))
				case 3:
					article.ReplyCount, _ = strconv.Atoi(strings.TrimSpace(s.Text()))
				case 4:
					val := strings.TrimSpace(s.Text())
					article.ReplyAt, _ = time.Parse("01-02 15:04", val)
				}
			})
			c.articleMap[article.Title] = article
		}
	})

	//get prev & next page
	doc.Find(".short-pages-2 a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if s.Text() == "上一页" {
			c.prevPage = href
		} else if s.Text() == "下一页" {
			c.nextPage = href
		}
	})
	return nil
}

func (c *Column) ListArticles() []*Article {
	articles := make([]*Article, 0, len(c.articleMap))
	for _, article := range c.articleMap {
		articles = append(articles, article)
	}
	return articles
}

func NewColumn(id, name, href string, client *resty.Client) *Column {
	return &Column{
		id:         id,
		name:       name,
		href:       href,
		articleMap: map[string]*Article{},
	}
}
