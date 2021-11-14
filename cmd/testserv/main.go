package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
)

func restHandler(w http.ResponseWriter, r *http.Request) {
	var delay float64 = 1.5 + rand.Float64()*1
	log.Infof("delaying for %f seconds", delay)
	time.Sleep(time.Duration(delay) * time.Second)
	fmt.Fprintf(w, "done, delay %f seconds\n", delay)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", restHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", "1235"),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err := srv.ListenAndServe()
	log.Fatal(err)
}
