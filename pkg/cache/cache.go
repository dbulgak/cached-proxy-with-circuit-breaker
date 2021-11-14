package cache

const Nil = CacheError("nil")

type CacheError string

func (e CacheError) Error() string { return string(e) }

type Cache interface {
	Get(key string) (value []byte, err error)
	Set(key string, value []byte) (err error)
}
