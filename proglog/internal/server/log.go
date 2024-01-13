package server

import (
	"fmt"
	"sync"
)

// - コミットログは時間順のレコードの並びで、追加だけ可能なデータ構造である
// - スライスを使って単純なコミットログを実装できる
type Log struct {
	mu      sync.Mutex
	records []Record
}

func NewLog() *Log {
	return &Log{}
}

// レコードをログに追加するメソッド
func (c *Log) Append(record Record) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	record.Offset = uint64(len(c.records))
	// recordsスライスにレコードを追加する
	c.records = append(c.records, record)

	return record.Offset, nil
}

// スライス内のレコードを取得するメソッド
func (c *Log) Read(offset uint64) (Record, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if offset >= uint64(len(c.records)) {
		return Record{}, ErrOffsetNotFound
	}

	return c.records[offset], nil
}

type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

var ErrOffsetNotFound = fmt.Errorf("offset not found")
