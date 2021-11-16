package proxy

import (
	"cachedproxy/pkg/cache"
	"cachedproxy/pkg/data"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sony/gobreaker"
	"io/ioutil"
	"net/http"
	"strings"
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
		cb:       cb,
		client:   client,
		settings: proxySettings,
	}

	return proxy, nil
}

func (p *Proxy) Request(username string, password string, req *data.DecodedRequest) (response []byte, isCached bool, err error) {
	log.WithField("req", req).Infof("REQUEST")

	val, err := p.client.Get(req)
	if err == nil {
		log.Infof("HIT: got key from cache client")
		return []byte(val), true, nil
	} else if err != nil && err != cache.Nil {
		log.Errorf("cache client get value error: %s, skipping", err)
	}

	log.Infof("MISS: no key in cache client")

	body, err := p.cb.Execute(func() (interface{}, error) {
		client := &http.Client{
			Timeout: p.settings.HttpTimeoutMs,
		}

		method := strings.ToUpper(req.Method)
		if method != "POST" && method != "GET" {
			return nil, fmt.Errorf("unexpected method %s", method)
		}

		req, err := http.NewRequest(req.Method, req.Url, strings.NewReader(req.Body))
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
		return nil, false, fmt.Errorf("got error %s", err.Error())
	}

	log.WithField("body", string(body.([]byte))).Infof("SAVE: saving response to client")
	err = p.client.Set(req, body.([]byte))
	if err != nil {
		log.Errorf("cache client set value error: %s, skipping", err)
	}

	return body.([]byte), false, nil
}
