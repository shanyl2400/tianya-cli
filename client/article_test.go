package client

import "testing"

func TestNextPage(t *testing.T) {
	article := NewArticle()
	article.Href = "/post-no05-506542-23.shtml"

	err := article.Open()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 500; i++ {
		content, err := article.Prev()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(content)
	}

}
