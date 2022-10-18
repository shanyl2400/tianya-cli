package cui

import (
	"os"
	"shanyl2400/tianya/client"

	"github.com/eiannone/keyboard"
	"github.com/muesli/termenv"
)

type CUI struct {
	client *client.Client
	term   *termenv.Output

	column  *client.Column
	article *client.Article
}

func (c *CUI) Start() {
	c.term = termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.TrueColor))
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}

	c.HomePage()
	c.Columns()
}

func NewCUI(c *client.Client) *CUI {
	return &CUI{
		client: c,
	}
}
