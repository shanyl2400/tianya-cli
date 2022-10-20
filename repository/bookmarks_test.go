package repository

import (
	"os"
	"shanyl2400/tianya/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	bookmark1 = &model.Bookmark{
		Title:      "title1",
		Author:     "author1",
		ViewCount:  1,
		ReplyCount: 1,

		// ReplyAt: time.Now(),

		Href: "href1",
		Type: "type0",

		Posts: []string{"post1", "post2"},
		Index: 1,

		Content:        "content1",
		CurViewContent: "viewcontent",
		Offset:         1,

		PrevPage: "prevpage1",
		NextPage: "nextpage1",
	}
	bookmark2 = &model.Bookmark{
		Title:      "title2",
		Author:     "author2",
		ViewCount:  2,
		ReplyCount: 2,

		// ReplyAt: time.Now(),

		Href: "href2",
		Type: "type2",

		Posts: []string{"post3", "post4"},
		Index: 2,

		Content:        "content2",
		CurViewContent: "viewcontent2",
		Offset:         2,

		PrevPage: "prevpage2",
		NextPage: "nextpage2",
	}
	bookmark3 = &model.Bookmark{
		Title:      "title3",
		Author:     "author3",
		ViewCount:  3,
		ReplyCount: 3,

		// ReplyAt: time.Now(),

		Href: "href3",
		Type: "type3",

		Posts: []string{"post5", "post4"},
		Index: 3,

		Content:        "content3",
		CurViewContent: "viewcontent3",
		Offset:         3,

		PrevPage: "prevpage3",
		NextPage: "nextpage3",
	}
	bookmark4 = &model.Bookmark{
		Title:      "title4",
		Author:     "author4",
		ViewCount:  4,
		ReplyCount: 4,

		// ReplyAt: time.Now(),

		Href: "href4",
		Type: "type4",

		Posts: []string{"post6", "post7"},
		Index: 4,

		Content:        "content4",
		CurViewContent: "viewcontent4",
		Offset:         4,

		PrevPage: "prevpage4",
		NextPage: "nextpage4",
	}

	bookmarkMap = map[string]*model.Bookmark{
		bookmark1.Title: bookmark1,
		bookmark2.Title: bookmark2,
		bookmark3.Title: bookmark3,
		bookmark4.Title: bookmark4,
	}
)

func TestSaveBookmarks(t *testing.T) {
	bm := GetBookmark()

	t.Log("opening data file...")
	err := bm.Open("./data")
	assert.NoError(t, err)

	defer func() {
		err = bm.Close()
		assert.NoError(t, err)

		os.RemoveAll("./data")
	}()

	t.Log("saving bookmarks...")
	for _, v := range bookmarkMap {
		err = bm.Save(v)
		assert.NoError(t, err)
	}

	t.Log("checking bookmarks load...")
	for k, v := range bookmarkMap {
		b, err := bm.Load(k)
		assert.NoError(t, err)
		assert.Equal(t, v, b)
	}

	t.Log("checking bookmarks list...")
	bs, err := bm.List()
	assert.NoError(t, err)
	assert.Equal(t, 4, len(bs))
	for i := range bs {
		b := bookmarkMap[bs[i].Title]
		assert.Equal(t, b, bs[i])
	}

	t.Log("Done")
}

func TestDeleteBookmarks(t *testing.T) {
	bm := GetBookmark()

	t.Log("opening data file...")
	err := bm.Open("./data")
	assert.NoError(t, err)

	defer func() {
		err = bm.Close()
		assert.NoError(t, err)

		os.RemoveAll("./data")
	}()

	t.Log("saving bookmarks...")
	for _, v := range bookmarkMap {
		err = bm.Save(v)
		assert.NoError(t, err)
	}

	t.Log("deleting bookmark1...")
	err = bm.Delete(bookmark1.Title)
	assert.NoError(t, err)

	t.Log("checking bookmarks load...")
	for k, v := range bookmarkMap {
		b, err := bm.Load(k)
		if k == bookmark1.Title {
			assert.Error(t, err)
			assert.Nil(t, b)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, v, b)
		}
	}

	t.Log("checking bookmarks list...")
	bs, err := bm.List()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(bs))
	for i := range bs {
		b := bookmarkMap[bs[i].Title]
		assert.Equal(t, b, bs[i])
	}

	t.Log("Done")
}

func TestUpdateBookmarks(t *testing.T) {
	bm := GetBookmark()

	t.Log("opening data file...")
	err := bm.Open("./data")
	assert.NoError(t, err)

	defer func() {
		err = bm.Close()
		assert.NoError(t, err)

		os.RemoveAll("./data")
	}()

	t.Log("saving bookmarks...")
	for _, v := range bookmarkMap {
		err = bm.Save(v)
		assert.NoError(t, err)
	}

	t.Log("updaing bookmark1...")
	bookmark1.NextPage = "foo"
	bookmark1.PrevPage = "bar"
	err = bm.Save(bookmark1)
	assert.NoError(t, err)

	t.Log("checking bookmark1...")
	b, err := bm.Load(bookmark1.Title)
	assert.NoError(t, err)
	assert.Equal(t, bookmark1, b)
	assert.Equal(t, "foo", b.NextPage)
	assert.Equal(t, "bar", b.PrevPage)

	t.Log("Done")
}
