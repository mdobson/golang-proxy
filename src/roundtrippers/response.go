package roundtrippers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

type ResponseInterceptTransport struct {
	http.RoundTripper
}

func (t *ResponseInterceptTransport) RoundTrip(r *http.Request) (resp *http.Response, err error) {
	log.WithFields(log.Fields{
		"middleware": "response-roundtripper",
	}).Info("Receieving response from target")
	res, err := t.RoundTripper.RoundTrip(r)

	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	b = bytes.Replace(b, []byte("Guest"), []byte("you clown!"), -1)
	body := ioutil.NopCloser(bytes.NewReader(b))
	res.Body = body
	res.ContentLength = int64(len(b))
	res.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return res, nil
}
