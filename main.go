package main

import (
	"dc-cb/cache"
	"dc-cb/proxy"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"net/http"
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

type Request struct {
	Url string
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

func (app *App) restHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)

	if ok {
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := app.Proxy.Request(username, password, req.Url)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		w.Write(resp)
	} else {
		// io.WriteString(w, "No basic auth, url is " + req.Url)
		fmt.Fprintf(w, "url %s, no basic auth", req.Url)
	}
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	settings := &Settings{
		RedisAddr:         os.Getenv("REDIS_ADDR"),
		MemcachedAddr:     os.Getenv("MEMCACHED_ADDR"),
		CacheType:         os.Getenv("CACHE_TYPE"),
		CacheExpirationMs: time.Duration(getIntEnv("CACHE_EXPIRATION_MS", 10000)) * time.Millisecond,
		HttpTimeoutMs:     time.Duration(getIntEnv("HTTP_TIMEOUT_MS", 2000)) * time.Millisecond,
		CbTimeoutMs:       time.Duration(getIntEnv("CB_TIMEOUT_MS", 60000)) * time.Millisecond,
	}
	log.Infoln(settings)

	app, err := NewApp(settings)
	if err != nil {
		log.Fatal(err)
	}

	listenPort := getEnv("LISTEN_PORT", "4000")

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.restHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", listenPort),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
