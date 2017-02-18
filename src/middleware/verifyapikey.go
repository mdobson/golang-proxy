package middleware

import (
	"data"
	"encoding/json"
	"net/http"
	"net/url"
)

type APIKeyContext struct {
	Result APIKeyVerifyResult `json:"result"`
	Type   string             `json:"string"`
}

type APIKeyVerifyResult struct {
	Key              string `json:"key"`
	ExpiredAt        int    `json:"expiresAt"`
	IssuedAt         int    `json:"issuedAt"`
	Status           string `json:"status"`
	RedirectionURIs  string `json:"redirectionURIs"`
	DeveloperID      string `json:"developerId"`
	DeveloperAppName string `json:"developerAppName"`
	ErrorCode        string `json:"errorCode"`
	Reason           string `json:"reasons"`
}

type VerifyApiKey struct {
	Proxy data.ProxyData
}

func (v VerifyApiKey) GetID() string {
	return "VerifyApiKey"
}

func (v VerifyApiKey) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-key")

		if apiKey != "" {
			form := url.Values{}
			form.Add("scopeuuid", v.Proxy.ProxyScope)
			form.Add("key", apiKey)
			form.Add("uriPath", v.Proxy.Basepath)
			form.Add("action", "verify")

			apiUrl := "http://localhost:9090/verifiers/apikey"
			resp, err := http.PostForm(apiUrl, form)

			if err != nil {
				http.Error(w, http.StatusText(500), 500)
			}

			var apiKeyContext APIKeyContext
			decoder := json.NewDecoder(resp.Body)
			decoder.UseNumber()

			if err := decoder.Decode(&apiKeyContext); err != nil {
				http.Error(w, http.StatusText(500), 500)
			}

			if verifyError := apiKeyContext.Type; verifyError == "ErrorResult" {
				http.Error(w, http.StatusText(401), 401)
			} else {
				status := apiKeyContext.Result.Status
				if status == "APPROVED" {
					next.ServeHTTP(w, r)
				}
			}

		} else {
			http.Error(w, http.StatusText(401), 401)
		}
		//next.ServeHTTP(w, r)
	})
}
