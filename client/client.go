package client

import (
	"net/http"
	"regexp"
	"shanyl2400/tianya/log"
	"shanyl2400/tianya/repository"

	"github.com/go-resty/resty/v2"
)

const (
	baseURL = "http://bbs.tianya.cn"

	menuPattern = "<li><a itemid=\"(.+)\" class=\"child_link\" href=\"(.+)\">(.+)</a></li>"
)

var (
	filteredColumn = map[string]bool{
		"天涯杂谈": true,
		"煮酒论史": true,
		"国际观察": true,
		"科技论坛": true,
		"仗剑天涯": true,
	}
)

type Client struct {
	columns    map[string]*Column
	httpClient *resty.Client
}

func (m *Client) Open() error {
	resp, err := m.access()
	if err != nil {
		log.WithField("err", err).
			Error("Access website failed")
		return err
	}
	r, err := regexp.Compile(menuPattern)
	if err != nil {
		log.WithField("err", err).
			WithField("pattern", menuPattern).
			Error("Compile menu failed")
		return err
	}

	groups := r.FindAllStringSubmatch(resp, -1)
	for i := range groups {
		if len(groups[i]) != 4 {
			log.WithField("groups", groups).
				WithField("resp", resp).
				Error("Invalid menu format")
			return err
		}
		column := NewColumn(groups[i][1], groups[i][3], groups[i][2], m.httpClient)
		m.columns[column.name] = column
	}

	//open bookmark
	err = repository.GetBookmark().Open("./bookmark")
	if err != nil {
		log.WithField("err", err).
			Error("Open bookmark failed")
		return err
	}

	return nil
}

func (m *Client) Close() {
	err := repository.GetBookmark().Close()
	if err != nil {
		log.WithField("err", err).
			Error("Close bookmark failed")
	}
}

func (m *Client) ListColumns() []*Column {
	columns := make([]*Column, 0, len(m.columns))
	for _, column := range m.columns {
		if m.filter(column.name) {
			columns = append(columns, column)
		}
	}
	return columns
}

func (m *Client) access() (string, error) {
	resp, err := m.httpClient.R().
		EnableTrace().
		Get(baseURL)
	if err != nil {
		log.WithField("err", err).
			WithField("url", baseURL).
			Error("access website failed")
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		log.WithField("err", err).
			WithField("statusCode", resp.StatusCode).
			WithField("res", resp).
			Error("access website failed")
		return "", ErrHTTPBadStateCode
	}
	return resp.String(), nil
}

func (m *Client) filter(name string) bool {
	return filteredColumn[name]
}

func NewClient() *Client {
	return &Client{
		httpClient: resty.New(),
		columns:    make(map[string]*Column),
	}
}
