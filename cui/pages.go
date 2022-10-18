package cui

import (
	"fmt"
	"os"
	"shanyl2400/tianya/client"

	"github.com/eiannone/keyboard"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

func (c *CUI) Columns() {
	spinnerInfo, err := pterm.DefaultSpinner.Start("Loading columns...")
	if err != nil {
		panic(err)
	}
	columns := c.client.ListColumns()
	spinnerInfo.Success("Load columns success")

	var options []string
	for i := range columns {
		options = append(options, fmt.Sprintf("%v", columns[i].Name()))
	}
	options = append(options, "> 退出")

	interactiveSelect := pterm.DefaultInteractiveSelect
	interactiveSelect.MaxHeight = 5
	selectedOption, _ := interactiveSelect.WithOptions(options).Show("Please select a column")
	pterm.Info.Printfln("Selected option: %s", pterm.Green(selectedOption))

	if selectedOption == "> 退出" {
		os.Exit(0)
	}

	for i := range columns {
		if columns[i].Name() == selectedOption {
			c.column = columns[i]
			c.Articles()
		}
	}
}

func (c *CUI) Articles() {
	c.term.ClearScreen()

	spinnerInfo, err := pterm.DefaultSpinner.Start("Loading articles...")
	if err != nil {
		panic(err)
	}

	err = c.column.Open()
	if err != nil {
		panic(err)
	}
	spinnerInfo.Success("Load articles success")

	articles := c.column.ListArticles()
	options := []string{"< 返回"}
	for i := range articles {
		options = append(options, fmt.Sprintf("%v", articles[i].Title))
	}

	interactiveSelect := pterm.DefaultInteractiveSelect
	interactiveSelect.MaxHeight = 5
	selectedOption, _ := interactiveSelect.WithOptions(options).Show("Please select an article")
	pterm.Info.Printfln("Selected option: %s", pterm.Green(selectedOption))

	if selectedOption == "< 返回" {
		c.Columns()
	} else {
		for i := range articles {
			if articles[i].Title == selectedOption {
				c.article = articles[i]
				c.Article()
			}
		}
	}
}

func (c *CUI) Article() {
	c.term.ClearScreen()
	spinnerInfo, err := pterm.DefaultSpinner.Start("Loading article...")
	if err != nil {
		panic(err)
	}

	err = c.article.Open()
	if err != nil {
		panic(err)
	}
	spinnerInfo.Success("Load article success")

	fmt.Print("Press d to start your read:")
	c.read(c.article)
}

func (c *CUI) HomePage() {
	c.term.ClearScreen()
	pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("T", pterm.NewStyle(pterm.FgCyan)),
		putils.LettersFromStringWithStyle("IANYA", pterm.NewStyle(pterm.FgLightMagenta))).
		Render()

	spinnerInfo, err := pterm.DefaultSpinner.Start("Connecting to the server...")
	if err != nil {
		panic(err)
	}
	err = c.client.Open()
	if err != nil {
		panic(err)
	}
	spinnerInfo.Success("Connected to server")
}

func (c *CUI) read(a *client.Article) {
out:
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		switch key {
		case keyboard.KeyEsc:
			break out
		case keyboard.KeyArrowLeft:
			c.prevPage(a)
		case keyboard.KeyArrowRight:
			c.nextPage(a)
		}

		switch char {
		case 'a':
			c.prevPage(a)
		case 'd':
			c.nextPage(a)
		case 'q':
			break out
		}
	}
	c.Articles()
}

func (c *CUI) prevPage(a *client.Article) {
	post, err := a.Prev()
	if err == client.ErrPageOutRange || err == client.ErrNoHistory {
		return
	}
	if err != nil {
		panic(err)
	}
	c.term.ClearScreen()
	pterm.DefaultParagraph.WithMaxWidth(60).Println(post)
}

func (c *CUI) nextPage(a *client.Article) {
	content, err := a.Next()
	if err == client.ErrPageOutRange {
		return
	}
	if err != nil {
		panic(err)
	}
	c.term.ClearScreen()
	pterm.DefaultParagraph.WithMaxWidth(60).Println(content)
}
