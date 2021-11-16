package app

import (
	"cachedproxy/pkg/cache"
	"cachedproxy/pkg/data"
	"cachedproxy/pkg/proxy"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type App struct {
	Proxy *proxy.Proxy
}

type Settings struct {
	RedisAddr         string
	MemcachedAddr     string
	CacheType         string
	CacheExpirationMs time.Duration
	CbTimeoutMs       time.Duration
	HttpTimeoutMs     time.Duration
}

func NewApp(settings *Settings) (*App, error) {
	client, err := getCacheClient(settings)
	if err != nil {
		return nil, err
	}

	prx, _ := proxy.NewProxy(client, &proxy.Settings{
		CbTimeoutMs:   settings.CbTimeoutMs,
		HttpTimeoutMs: settings.HttpTimeoutMs,
	})
	app := &App{
		Proxy: prx,
	}

	return app, nil
}

func getCacheClient(settings *Settings) (cache.Cache, error) {
	var client cache.Cache

	log.Infof("%s cache type", settings.CacheType)

	switch settings.CacheType {
	case "redis":
		if strings.Compare(settings.RedisAddr, "") == 0 {
			return nil, fmt.Errorf("no REDIS_ADDR setting")
		}
		client = cache.NewRedis(&cache.RedisSettings{
			Addr:       settings.RedisAddr,
			Password:   "",
			DB:         0,
			Expiration: settings.CacheExpirationMs,
		})
	case "memcached":
		if strings.Compare(settings.MemcachedAddr, "") == 0 {
			return nil, fmt.Errorf("no MEMCACHED_ADDR setting")
		}
		client = cache.NewMemcached(&cache.MemcachedSettings{
			Url:        settings.MemcachedAddr,
			Expiration: settings.CacheExpirationMs,
		})
	default:
		return nil, fmt.Errorf("unexpected cache type %s", settings.CacheType)
	}

	return client, nil
}

func (app *App) RestHandler(w http.ResponseWriter, r *http.Request) {
	var req data.Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	decodedRequest, err := data.DecodeRequest(req)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username, password, _ := r.BasicAuth()
	resp, isCached, err := app.Proxy.Request(username, password, decodedRequest)
	w.Header().Set("X-Cache", strconv.FormatBool(isCached))
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Write(resp)
}
