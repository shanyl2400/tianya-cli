package repository

import (
	"encoding/json"
	"errors"
	"shanyl2400/tianya/model"

	"go.etcd.io/bbolt"
)

var (
	ErrNoSuchBookmark = errors.New("no such bookmark")
)

type Bookmarks interface {
	Open(path string) error
	Close() error

	Save(b *model.Bookmark) error
	Load(title string) (*model.Bookmark, error)
	Delete(title string) error

	List() ([]*model.Bookmark, error)
}

type BoltDBBookmarks struct {
	db   *bbolt.DB
	path string
}

func (b *BoltDBBookmarks) Save(bookmark *model.Bookmark) error {
	err := b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(b.bucket())

		buf, err := json.Marshal(bookmark)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket.
		return bucket.Put([]byte(bookmark.Title), buf)
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BoltDBBookmarks) Delete(title string) error {
	err := b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(b.bucket())

		// Persist bytes to users bucket.
		return bucket.Delete([]byte(title))
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BoltDBBookmarks) Load(title string) (*model.Bookmark, error) {
	var bookmark *model.Bookmark
	err := b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(b.bucket())

		raw := bucket.Get([]byte(title))
		if raw == nil {
			return ErrNoSuchBookmark
		}
		bookmark = new(model.Bookmark)
		err := json.Unmarshal(raw, bookmark)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return bookmark, nil
}

func (b *BoltDBBookmarks) List() ([]*model.Bookmark, error) {
	bookmarks := make([]*model.Bookmark, 0)
	err := b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(b.bucket())
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			bookmark := new(model.Bookmark)
			err := json.Unmarshal(v, bookmark)
			if err != nil {
				return err
			}
			bookmarks = append(bookmarks, bookmark)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return bookmarks, nil
}

func (b *BoltDBBookmarks) Open(path string) error {
	b.path = path
	db, err := bbolt.Open(b.path, 0666, nil)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(b.bucket())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	b.db = db
	return nil
}

func (b *BoltDBBookmarks) Close() error {
	return b.db.Close()
}

func (b *BoltDBBookmarks) bucket() []byte {
	return []byte("bookmarks")
}

var (
	_bookmark = new(BoltDBBookmarks)
)

func GetBookmark() Bookmarks {
	return _bookmark
}
