package cache

import (
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"time"
)

type RedisL struct {
	BaseCache
	client     *redis.Client
	expiration time.Duration
}

type RedisSettings struct {
	Addr       string
	Password   string
	DB         int
	Expiration time.Duration
}

func NewRedis(settings *RedisSettings) *RedisL {
	client := &RedisL{
		client: redis.NewClient(&redis.Options{
			Addr:     settings.Addr,
			Password: settings.Password,
			DB:       settings.DB,
		}),
		expiration: settings.Expiration,
	}
	return client
}

func (r *RedisL) Get(key string) ([]byte, error) {
	value, err := r.client.Get(r.GetHashedKey(key)).Result()
	if err == redis.Nil {
		return nil, Nil
	} else if err != nil {
		log.Errorf("redis error: %s", err)
		return nil, err
	}
	return []byte(value), err
}

func (r *RedisL) Set(key string, value []byte) (err error) {
	err = r.client.Set(r.GetHashedKey(key), value, r.expiration).Err()
	return err
}
