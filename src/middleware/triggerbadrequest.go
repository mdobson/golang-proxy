package middleware

import "net/http"

type TriggerBadRequest struct {
}

func (h *TriggerBadRequest) GetID() string {
	return "TriggerBadRequest"
}

func (h *TriggerBadRequest) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/error" {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		next.ServeHTTP(w, r)
	})
}
