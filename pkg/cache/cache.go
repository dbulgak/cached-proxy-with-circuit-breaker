package cache

import (
	"cachedproxy/pkg/data"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

const Nil = CacheError("nil")

type CacheError string

func (e CacheError) Error() string { return string(e) }

type Cache interface {
	Get(key string) (value []byte, err error)
	Set(key string, value []byte) (err error)
}

func GetHashedKey(req *data.DecodedRequest) string {
	hash := md5.Sum([]byte(req.String()))
	hashedkey := fmt.Sprintf("%s_%s", "cachedproxy", hex.EncodeToString(hash[:]))
	return hashedkey
}
