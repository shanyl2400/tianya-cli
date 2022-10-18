package client

import (
	"fmt"
	"net/http"
	"regexp"

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
		return err
	}
	r, err := regexp.Compile(menuPattern)
	if err != nil {
		return err
	}

	groups := r.FindAllStringSubmatch(resp, -1)
	for i := range groups {
		if len(groups[i]) != 4 {
			return err
		}
		column := NewColumn(groups[i][1], groups[i][3], groups[i][2], m.httpClient)
		m.columns[column.name] = column
	}

	return nil
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
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("access website failed, status: %v", resp.StatusCode())
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
