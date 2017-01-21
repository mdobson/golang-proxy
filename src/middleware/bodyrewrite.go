package middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type BodyRewrite struct {
}

func (h *BodyRewrite) GetID() string {
	return "BodyRewrite"
}

func (h *BodyRewrite) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			fmt.Println("Not a GET request. Let's rewrite the body.")
			newBodyContent := "Body Rewritten By Proxy!"
			r.Body = ioutil.NopCloser(strings.NewReader(newBodyContent))
			r.ContentLength = int64(len(newBodyContent))
		}
		next.ServeHTTP(w, r)
	})
}
