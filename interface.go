package marshalcache

type Cache interface {
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
}
type ObjectCache interface {
	Put(key string, value interface{}) error
	Get(key string) (interface{}, error)
}
type MarshalCache interface {
	Marshal(obj interface{}) ([]byte, error)
	Unmarshal(buf []byte, obj interface{}) error
}
type McObject interface {
	GetMcKey() string
}
type Marshaller func(obj interface{}) ([]byte, error)
type Unmarshaler func(buf []byte, obj interface{}) error
