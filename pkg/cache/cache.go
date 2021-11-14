package cache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
)

const Nil = CacheError("nil")

type CacheError string

func (e CacheError) Error() string { return string(e) }

type Cache interface {
	Get(key string) (value []byte, err error)
	Set(key string, value []byte) (err error)
}

type BaseCache struct {
}

func (b *BaseCache) GetHashedKey(key string) string {
	hash := md5.Sum([]byte(key))
	hashedkey := fmt.Sprintf("%s_%s", "cachedproxy", hex.EncodeToString(hash[:]))

	log.Infof("%s key converted to %s hash", key, hashedkey)

	return hashedkey
}
