package marshalcache

import (
	"encoding/json"
	"errors"
	"testing"
)

type memCache2 struct {
	cache map[string]interface{}
}

func newMemCache2() *memCache2 {
	return &memCache2{cache: make(map[string]interface{})}
}
func (m *memCache2) Put(key string, value interface{}) error {
	m.cache[key] = value
	return nil
}
func (m *memCache2) Get(key string) (interface{}, error) {
	v, ok := m.cache[key]
	if ok {
		return v, nil
	}
	return nil, errors.New("not found")
}

func BenchmarkMarshalCacheImpl2_Marshal_Direct(b *testing.B) {
	block := initBlock(1000)
	for i := 0; i < b.N; i++ {
		json.Marshal(block)
	}
}
func BenchmarkMarshalCacheImpl2_Marshal_UseCache(b *testing.B) {
	block := initBlock(1000)
	mcache := NewMarshalCache2(newMemCache2(), json.Marshal, json.Unmarshal, false)
	for i := 0; i < b.N; i++ {
		mcache.Marshal(block)
	}
}
func BenchmarkMarshalCacheImpl2_Marshal_NoCache(b *testing.B) {
	block := initBlock2(1000)
	mcache := NewMarshalCache2(newMemCache2(), json.Marshal, json.Unmarshal, false)
	for i := 0; i < b.N; i++ {
		mcache.Marshal(block)
	}
}
func BenchmarkMarshalCacheImpl2_Unmarshal_UseCache(b *testing.B) {
	data,_:=json.Marshal( initBlock(1000))
	mcache := NewMarshalCache2(newMemCache2(), json.Marshal, json.Unmarshal, false)
	for i := 0; i < b.N; i++ {
		var block *Block
		err:=mcache.Unmarshal(data,&block)
		if err!=nil{
			b.Fail()
		}
		mcache.Marshal(block)
	}
}
func BenchmarkMarshalCacheImpl2_Unmarshal_NoCache(b *testing.B) {
	data,_:=json.Marshal( initBlock2(1000))
	mcache := NewMarshalCache2(newMemCache2(), json.Marshal, json.Unmarshal, false)
	for i := 0; i < b.N; i++ {
		var block *Block2
		mcache.Unmarshal(data,&block)
		mcache.Marshal(block)
	}
}
func BenchmarkMarshalCacheImpl2_Unmarshal_Direct(b *testing.B) {
	data,_:=json.Marshal( initBlock(1000))
	for i := 0; i < b.N; i++ {
		var block *Block
		json.Unmarshal(data,&block)
		json.Marshal(block)
	}
}