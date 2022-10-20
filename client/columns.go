package client

import (
	"net/http"
	"shanyl2400/tianya/log"
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
	url  string

	articleMap map[string]*Article
	nextPage   string
	// prevPage   string
	prevPageList []string
}

func (c *Column) Open() error {
	return c.enterPage(c.href)
}

func (c *Column) Name() string {
	return c.name
}

func (c *Column) NextPage() error {
	if c.nextPage == "" {
		log.Info("no next page")
		return ErrPageOutRange
	}
	c.prevPageList = append(c.prevPageList, c.url)
	return c.enterPage(c.nextPage)
}

func (c *Column) PrevPage() error {
	if len(c.prevPageList) < 1 {
		log.Info("no prev page")
	}

	page := c.prevPageList[len(c.prevPageList)-1]
	c.prevPageList = c.prevPageList[:len(c.prevPageList)-1]

	return c.enterPage(page)
}

func (c *Column) enterPage(path string) error {
	c.url = path
	res, err := http.Get(baseURL + path)
	if err != nil {
		log.WithField("err", err).
			WithField("url", baseURL+path).
			Error("access website failed")
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.WithField("err", err).
			WithField("statusCode", res.StatusCode).
			WithField("res", res).
			Error("access website failed")
		return ErrHTTPBadStateCode
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.WithField("err", err).
			WithField("body", res.Body).
			WithField("res", res).
			Error("create dom failed")
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
		if s.Text() == "下一页" {
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

func (c *Column) HasNextPage() bool {
	return c.nextPage != ""
}

func (c *Column) HasPrevPage() bool {
	return len(c.prevPageList) > 0
}

func NewColumn(id, name, href string, client *resty.Client) *Column {
	return &Column{
		id:         id,
		name:       name,
		href:       href,
		articleMap: map[string]*Article{},
	}
}
