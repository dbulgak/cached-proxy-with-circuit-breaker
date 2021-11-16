package main

import (
	"cachedproxy/pkg/app"
	"cachedproxy/pkg/utils"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	settings := &app.Settings{
		RedisAddr:         os.Getenv("REDIS_ADDR"),
		MemcachedAddr:     os.Getenv("MEMCACHED_ADDR"),
		CacheType:         os.Getenv("CACHE_TYPE"),
		CacheExpirationMs: time.Duration(utils.GetIntEnv("CACHE_EXPIRATION_MS", 10000)) * time.Millisecond,
		HttpTimeoutMs:     time.Duration(utils.GetIntEnv("HTTP_TIMEOUT_MS", 2000)) * time.Millisecond,
		CbTimeoutMs:       time.Duration(utils.GetIntEnv("CB_TIMEOUT_MS", 60000)) * time.Millisecond,
	}
	log.Infoln(settings)

	application, err := app.NewApp(settings)
	if err != nil {
		log.Fatal(err)
	}

	listenPort := utils.GetEnv("LISTEN_PORT", "4000")

	mux := http.NewServeMux()
	mux.HandleFunc("/", application.RestHandler)

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
