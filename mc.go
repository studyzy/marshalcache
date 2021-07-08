package marshalcache

import "sync"

type marshalCacheImpl struct {
	cache Cache
	sync.Mutex
	mfn            Marshaller
	umfn           Unmarshaler
	openCacheError bool
}

func NewMarshalCache(cache Cache, marshalFunc Marshaller, unmarshalFunc Unmarshaler, openCacheError bool) MarshalCache {
	return &marshalCacheImpl{
		cache:          cache,
		mfn:            marshalFunc,
		umfn:           unmarshalFunc,
		openCacheError: openCacheError,
	}
}
func (mc *marshalCacheImpl) Marshal(obj interface{}) ([]byte, error) {
	mc.Lock()
	defer mc.Unlock()
	mcobj, ok := obj.(McObject)
	if !ok {
		return mc.mfn(obj)
	}
	key := mcobj.GetMcKey()
	cacheData, err := mc.cache.Get(key)
	if len(cacheData) > 0 { //found in cache
		return cacheData, err
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
func (mc *marshalCacheImpl) Unmarshal(buf []byte, obj interface{}) error {
	mc.Lock()
	defer mc.Unlock()
	err := mc.umfn(buf, obj)
	if err != nil {
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
