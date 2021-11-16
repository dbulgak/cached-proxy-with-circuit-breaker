package cache

import (
	"cachedproxy/pkg/data"
	"github.com/bradfitz/gomemcache/memcache"
	log "github.com/sirupsen/logrus"
	"time"
)

type MemcachedL struct {
	BaseCache
	client        *memcache.Client
	expirationSec int32
}

type MemcachedSettings struct {
	Url        string
	Expiration time.Duration
}

func NewMemcached(settings *MemcachedSettings) *MemcachedL {
	if settings.Expiration < time.Second {
		log.Fatal("cannot be less than 1 second")
	}

	memcached := &MemcachedL{
		client:        memcache.New(settings.Url),
		expirationSec: int32(settings.Expiration / time.Second),
	}
	return memcached
}

func (m *MemcachedL) Get(req *data.DecodedRequest) (value []byte, err error) {
	item, err := m.client.Get(m.GetHashedKey(req))
	if err == memcache.ErrCacheMiss {
		return nil, Nil
	} else if err != nil {
		return nil, err
	}

	return item.Value, err
}

func (m *MemcachedL) Set(req *data.DecodedRequest, value []byte) (err error) {
	err = m.client.Set(&memcache.Item{
		Key:        m.GetHashedKey(req),
		Value:      value,
		Flags:      0,
		Expiration: m.expirationSec,
	})
	return err
}
