package repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	items = []*HistoryItem{
		{
			Title:     "title1",
			Parititon: "paritition1",
			Content:   "content1",
		},
		{
			Title:     "title1",
			Parititon: "paritition1",
			Content:   "content2",
		},
		{
			Title:     "title1",
			Parititon: "paritition1",
			Content:   "content3",
		},
		{
			Title:     "title2",
			Parititon: "paritition2",
			Content:   "content4",
		},
		{
			Title:     "title2",
			Parititon: "paritition2",
			Content:   "content5",
		},
		{
			Title:     "title2",
			Parititon: "paritition2",
			Content:   "content6",
		},
	}
)

func TestPushPop(t *testing.T) {
	his := NewHistory("./data")
	err := his.Open()
	assert.NoError(t, err)

	for i := range items {
		id, err := his.Push(items[i].Title, items[i].Parititon, items[i].Content)
		items[i].ID = id
		assert.NoError(t, err)
	}

	for i := len(items) - 1; i >= 0; i-- {
		empty := his.IsEmpty(items[i].Title, items[i].Parititon)
		assert.False(t, empty)

		c, err := his.Pop(items[i].Title, items[i].Parititon)
		assert.NoError(t, err)
		assert.Equal(t, items[i], c)
		t.Logf("item: %#v", items[i])
	}
	err = his.Close()
	assert.NoError(t, err)

	os.RemoveAll("./data")
}
