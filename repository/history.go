package repository

import (
	"encoding/binary"
	"encoding/json"

	"go.etcd.io/bbolt"
)

type HistoryItem struct {
	ID        uint64
	Title     string
	Parititon string
	Content   string
}

type History interface {
	Open() error
	Close() error

	Push(title, partition, content string) (uint64, error)
	Pop(title, partition string) (*HistoryItem, error)

	IsEmpty(title, partition string) bool
}

type BoltDBHistory struct {
	db   *bbolt.DB
	path string
}

func (h *BoltDBHistory) Push(title, partition, content string) (uint64, error) {
	var id uint64
	err := h.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(h.bucket(title, partition))
		if err != nil {
			return err
		}
		id, _ = b.NextSequence()

		item := HistoryItem{
			ID:        id,
			Title:     title,
			Parititon: partition,
			Content:   content,
		}

		buf, err := json.Marshal(item)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket.
		return b.Put(itob(id), buf)
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (h *BoltDBHistory) Pop(title, partition string) (*HistoryItem, error) {
	item := new(HistoryItem)
	err := h.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(h.bucket(title, partition))
		if err != nil {
			return err
		}
		cursor := b.Cursor()

		key, value := cursor.Last()

		err = json.Unmarshal(value, item)
		if err != nil {
			return err
		}
		err = b.Delete(key)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (h *BoltDBHistory) IsEmpty(title, partition string) bool {
	flag := false
	h.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(h.bucket(title, partition))
		if err != nil {
			flag = true
			return err
		}
		cursor := b.Cursor()
		key, _ := cursor.First()
		if key == nil {
			flag = true
		}
		return nil
	})
	return flag
}
func (h *BoltDBHistory) Open() error {
	db, err := bbolt.Open(h.path, 0666, nil)
	if err != nil {
		return err
	}
	h.db = db
	return nil
}

func (h *BoltDBHistory) Close() error {
	return h.db.Close()
}

func (h BoltDBHistory) bucket(title, partition string) []byte {
	return []byte("history/" + title + "/" + partition)
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func NewHistory(path string) History {
	return &BoltDBHistory{
		path: path,
	}
}
