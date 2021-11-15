package cache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

const Nil = CacheError("nil")

type CacheError string

func (e CacheError) Error() string { return string(e) }

type Cache interface {
	Get(req Request) (value []byte, err error)
	Set(req Request, value []byte) (err error)
}

type BaseCache struct {
}

func (b *BaseCache) GetHashedKey(req Request) string {
	hash := md5.Sum([]byte(req.String()))
	hashedkey := fmt.Sprintf("%s_%s", "cachedproxy", hex.EncodeToString(hash[:]))
	return hashedkey
}

type Request struct {
	Url    string
	Method string
	Body   string
}

func (r Request) String() string {
	return fmt.Sprintf("%s_%s_%s", r.Url, r.Method, r.Body)
}
