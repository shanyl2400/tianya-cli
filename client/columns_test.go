package client

import (
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestAccessColumn(t *testing.T) {
	column := NewColumn("free", "天涯杂谈", "/list-free-1.shtml", resty.New())
	err := column.Open()
	if err != nil {
		panic(err)
	}
	articles := column.ListArticles()
	article := articles[0]

	article.Open()
	post, _ := article.NextPost()
	for post != nil {
		fmt.Println(post.Content)
		fmt.Println("------------")
		post, _ = article.NextPost()
	}

	// t.Log("----------next page----------")
	// err = column.NextPage()
	// if err != nil {
	// 	panic(err)
	// }
	// articles = column.ListArticles()
	// for i := range articles {
	// 	t.Logf("article title: %v, author: %v", articles[i].Title, articles[i].Author)
	// }

	// t.Log("----------prev page----------")
	// err = column.PrevPage()
	// if err != nil {
	// 	panic(err)
	// }
	// articles = column.ListArticles()
	// for i := range articles {
	// 	t.Logf("article title: %v, author: %v", articles[i].Title, articles[i].Author)
	// }
}
