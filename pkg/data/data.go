package data

import (
	"encoding/base64"
	"fmt"
)

type Request struct {
	Url    string
	Method string
	Body   string
}

func (r Request) String() string {
	return fmt.Sprintf("%s_%s_%s", r.Url, r.Method, r.Body)
}

type DecodedRequest struct {
	Url    string
	Method string
	Body   string
}

func (r DecodedRequest) String() string {
	return fmt.Sprintf("%s_%s_%s", r.Url, r.Method, r.Body)
}

func DecodeRequest(req Request) (*DecodedRequest, error) {
	decodedBody, err := base64.StdEncoding.DecodeString(req.Body)
	if err != nil {
		return nil, err
	}

	decodedRequest := &DecodedRequest{
		Url:    req.Url,
		Method: req.Method,
		Body:   string(decodedBody),
	}

	return decodedRequest, nil
}
