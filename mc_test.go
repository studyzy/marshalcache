package marshalcache

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"
)

type memCache struct {
	cache map[string][]byte
}

func newMemCache() *memCache {
	return &memCache{cache: make(map[string][]byte)}
}
func (m *memCache) Put(key string, value []byte) error {
	m.cache[key] = value
	return nil
}
func (m *memCache) Get(key string) ([]byte, error) {
	v, ok := m.cache[key]
	if ok {
		return v, nil
	}
	return nil, errors.New("not found")
}

type Block struct {
	BlockHash []byte
	Header    *Header
	Txs       []*Tx
}

func (b *Block) GetMcKey() string {
	return fmt.Sprintf("Block[%x]", b.BlockHash)
}

type Header struct {
	Height    uint64
	PreHash   []byte
	Nonce     uint64
	Diff      uint64
	Root      []byte
	Timestamp int64
}
type Tx struct {
	TxId      string
	From      string
	To        string
	Amount    uint64
	Signature []byte
}

func (tx *Tx) GetMcKey() string {
	return fmt.Sprintf("Tx[%s]", tx.TxId)
}

//Block2没有GetMcKey，无法用到缓存，用于对比性能
type Block2 struct {
	BlockHash []byte
	Header    *Header
	Txs       []*Tx
}

func initBlock(txCount int) *Block {
	b := &Block{
		BlockHash: []byte("hash1"),
		Header: &Header{
			Height:    123,
			PreHash:   []byte("hhhhhhhhh"),
			Nonce:     11110,
			Diff:      2324230,
			Root:      []byte("ssssssssss"),
			Timestamp: 0,
		},
		Txs: make([]*Tx, txCount),
	}
	for i := 0; i < txCount; i++ {
		b.Txs[i] = &Tx{
			TxId:      "Tx" + strconv.Itoa(i),
			From:      "Addr" + strconv.Itoa(i),
			To:        "Addr" + strconv.Itoa(i+1000000),
			Amount:    uint64(i) * 100,
			Signature: []byte("signature" + strconv.Itoa(i*1000)),
		}
	}
	return b
}
func initBlock2(txCount int) *Block2 {
	b := initBlock(txCount)
	return &Block2{
		BlockHash: b.BlockHash,
		Header:    b.Header,
		Txs:       b.Txs,
	}
}
func BenchmarkMarshalCacheImpl_Marshal_Direct(b *testing.B) {
	block := initBlock(1000)
	for i := 0; i < b.N; i++ {
		json.Marshal(block)
	}
}
func BenchmarkMarshalCacheImpl_Marshal_UseCache(b *testing.B) {
	block := initBlock(1000)
	mcache := NewMarshalCache(newMemCache(), json.Marshal, json.Unmarshal, false)
	for i := 0; i < b.N; i++ {
		mcache.Marshal(block)
	}
}
func BenchmarkMarshalCacheImpl_Marshal_NoCache(b *testing.B) {
	block := initBlock2(1000)
	mcache := NewMarshalCache(newMemCache(), json.Marshal, json.Unmarshal, false)
	for i := 0; i < b.N; i++ {
		mcache.Marshal(block)
	}
}
