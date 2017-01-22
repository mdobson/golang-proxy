package middleware

import (
	"net/http"

	"github.com/google/uuid"
)

type HeaderSet struct {
}

func (h HeaderSet) GetID() string {
	return "HeaderSet"
}

func (h HeaderSet) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gatewayFlowID, _ := uuid.NewUUID()

		//Sets a header in the outgoing request
		r.Header.Add("X-ID", gatewayFlowID.String())

		//Sets a header in the response
		w.Header().Set("X-GoProxy", "GoProxy")
		gatewayResponseID, _ := uuid.NewUUID()
		w.Header().Set("X-Gateway-Response-ID", gatewayResponseID.String())

		next.ServeHTTP(w, r)
	})
}
