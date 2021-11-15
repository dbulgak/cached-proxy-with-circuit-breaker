package cache

import (
	"cachedproxy/pkg/app"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

const Nil = CacheError("nil")

type CacheError string

func (e CacheError) Error() string { return string(e) }

type Cache interface {
	Get(req app.Request) (value []byte, err error)
	Set(req app.Request, value []byte) (err error)
}

type BaseCache struct {
}

func (b *BaseCache) GetHashedKey(req app.Request) string {
	hash := md5.Sum([]byte(req.String()))
	hashedkey := fmt.Sprintf("%s_%s", "cachedproxy", hex.EncodeToString(hash[:]))
	return hashedkey
}
