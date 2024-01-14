package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	// レコードサイズとインデックスエントリを永続化するためのエンコーディングを定義
	enc = binary.BigEndian
)

const (
	// レコードの長さを格納するために使うバイト数を定義
	lenWidth = 8
)

// ファイル操作をカプセル化する構造体
// ファイルへの書き込みや読み出しを行う
type store struct {
	*os.File
	mu   sync.Mutex    // 同時書き込みを防ぐためのミューテックス
	buf  *bufio.Writer // バッファリングを行うライター
	size uint64        // 現在のファイルサイズ
}

// 与えられたファイルに対する`store`を作成する関数
func newStore(f *os.File) (*store, error) {
	// ファイルの現在のサイズを取得
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(fi.Size())
	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil
}

// 新しいレコードをストアに追加するメソッド
// 追加されたレコードのサイズと位置を返す
func (s *store) Append(p []byte) (n uint64, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// ストアがファイル内でレコードを保持する位置を取得
	pos = s.size
	// レコードのサイズを先に書き込む
	if err := binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return 0, 0, err
	}
	w, err := s.buf.Write(p)
	if err != nil {
		return 0, 0, err
	}
	w += lenWidth
	s.size += uint64(w)
	return uint64(w), pos, nil
}

// 指定された位置に格納されているレコードを読み出すメソッド
func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}
	size := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}
	b := make([]byte, enc.Uint64(size))
	if _, err := s.File.ReadAt(b, int64(pos+lenWidth)); err != nil {
		return nil, err
	}
	return b, nil
}

// ストアのファイルのoffセットからpバイト読み出しするメソッド
func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.File.ReadAt(p, off)
}

func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.buf.Flush()
	if err != nil {
		return err
	}
	return s.File.Close()
}
