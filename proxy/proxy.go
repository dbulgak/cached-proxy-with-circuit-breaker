package proxy

import (
	"dc-cb/cache"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sony/gobreaker"
	"io/ioutil"
	"net/http"
	"time"
)

type Proxy struct {
	cb       *gobreaker.CircuitBreaker
	client   cache.Cache
	settings *Settings
}

type Settings struct {
	CbTimeoutMs   time.Duration
	HttpTimeoutMs time.Duration
}

func NewProxy(client cache.Cache, proxySettings *Settings) (*Proxy, error) {
	settings := gobreaker.Settings{
		Timeout: proxySettings.CbTimeoutMs,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	}
	cb := gobreaker.NewCircuitBreaker(settings)

	proxy := &Proxy{
		cb:     cb,
		client: client,
		settings: proxySettings,
	}

	return proxy, nil
}

func (p *Proxy) Request(username, password, url string) (response []byte, err error) {
	val, err := p.client.Get(url)
	if err == nil {
		log.Infof("HIT: got %s key, '%s' value from client", url, val)
		return []byte(val), nil
	} else if err != nil && err != cache.Nil {
		log.Errorf("cache client get value error: %s, skipping", err)
	}

	log.Infof("MISS: no %s key in client", url)

	body, err := p.cb.Execute(func() (interface{}, error) {
		client := &http.Client{
			Timeout: p.settings.HttpTimeoutMs,
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("got error %s", err.Error())
		}

		if username != "" || password != "" {
			req.SetBasicAuth(username, password)
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("got error %s", err.Error())
		}
		defer resp.Body.Close()

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return buf, fmt.Errorf("got error %s", err.Error())
		}

		return buf, nil
	})

	if err != nil {
		return nil, fmt.Errorf("got error %s", err.Error())
	}

	log.Infof("SAVE: saving %s response to client", url)
	err = p.client.Set(url, body.([]byte))
	if err != nil {
		log.Errorf("cache client set value error: %s, skipping", err)
	}

	return body.([]byte), nil
}
