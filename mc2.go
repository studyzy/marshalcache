package marshalcache

import (
	"fmt"
	"sync"
)

type marshalCacheImpl2 struct {
	cache ObjectCache
	sync.Mutex
	mfn            Marshaller
	umfn           Unmarshaler
	openCacheError bool
}

func NewMarshalCache2(cache ObjectCache, marshalFunc Marshaller, unmarshalFunc Unmarshaler, openCacheError bool) MarshalCache {
	return &marshalCacheImpl2{
		cache:          cache,
		mfn:            marshalFunc,
		umfn:           unmarshalFunc,
		openCacheError: openCacheError,
	}
}
func (mc *marshalCacheImpl2) Marshal(obj interface{}) ([]byte, error) {
	mc.Lock()
	defer mc.Unlock()
	mcobj, ok := obj.(McObject)
	if !ok {
		return mc.mfn(obj)
	}
	key := mcobj.GetMcKey()
	cacheData, err := mc.cache.Get(key)
	if cacheData!=nil { //found in cache
		bytes := cacheData.([]byte)
		return bytes, err

	}
	//not found in cache
	data, err := mc.mfn(obj)
	if err != nil {
		return data, err
	}
	err = mc.cache.Put(key, data)
	if err != nil && mc.openCacheError {
		return data, err
	}
	return data, nil
}
func (mc *marshalCacheImpl2) Unmarshal(buf []byte, obj interface{}) error {
	mc.Lock()
	defer mc.Unlock()
	hexKey:=fmt.Sprintf("%x",buf)
	cacheObj,err:= mc.cache.Get(hexKey)
	if err==nil &&cacheObj!=nil{
		 obj=cacheObj
		 return nil
	}
	//no cache hexKey
	err = mc.umfn(buf, obj)
	if err != nil {
		return err
	}
	err = mc.cache.Put(hexKey, obj)
	if err != nil && mc.openCacheError {
		return err
	}
	mcobj, ok := obj.(McObject)
	if !ok {
		return nil
	}
	//put buf in cache
	err = mc.cache.Put(mcobj.GetMcKey(), buf)
	if err != nil && mc.openCacheError {
		return err
	}
	return nil
}
