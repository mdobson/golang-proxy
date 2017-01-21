package middleware

import (
	"net/http"
)

type RequestMiddleware interface {
	GetID() string
	Handle(next http.Handler) http.Handler
}

// type RequestMiddlewareSequence struct {
// 	MiddlewareIDs []string
// }

// type MiddlewareRepository struct {
// 	Middlewares []RequestMiddleware
// }

// func (m *MiddlewareRepository) GetMiddleware(id string) RequestMiddleware {

// }
