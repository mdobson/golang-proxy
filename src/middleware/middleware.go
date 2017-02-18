package middleware

import "net/http"
import "data"

type RequestMiddleware interface {
	GetID() string
	Handle(next http.Handler) http.Handler
}

type RequestMiddlewareSequence struct {
	MiddlewareIDs map[string]RequestMiddleware
}

func CreateService(proxyData data.ProxyData) RequestMiddlewareSequence {
	sequence := RequestMiddlewareSequence{MiddlewareIDs: make(map[string]RequestMiddleware)}
	v := VerifyApiKey{
		Proxy: proxyData,
	}
	sequence.RegisterMiddleware(v.GetID(), v)
	h := HeaderSet{}
	sequence.RegisterMiddleware(h.GetID(), h)
	b := BodyRewrite{}
	sequence.RegisterMiddleware(b.GetID(), b)
	e := TriggerBadRequest{}
	sequence.RegisterMiddleware(e.GetID(), e)
	return sequence
}

func (r *RequestMiddlewareSequence) RegisterMiddleware(name string, ware RequestMiddleware) {
	r.MiddlewareIDs[name] = ware
}

func (r *RequestMiddlewareSequence) GetMiddlewareHandler(name string) RequestMiddleware {
	return r.MiddlewareIDs[name]
}

func (r *RequestMiddlewareSequence) Compile(middlewares []string, finalMiddleware http.Handler) http.Handler {
	var currentMiddleware = finalMiddleware
	for _, middlewareKey := range middlewares {
		middleware := r.GetMiddlewareHandler(middlewareKey)
		currentMiddleware = middleware.Handle(currentMiddleware)
	}
	return currentMiddleware
}

// type MiddlewareRepository struct {
// 	Middlewares []RequestMiddleware
// }

// func (m *MiddlewareRepository) GetMiddleware(id string) RequestMiddleware {

// }
