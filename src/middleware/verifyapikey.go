package middleware

import (
	"data"
	"encoding/json"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
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
		log.WithFields(log.Fields{
			"middleware": v.GetID(),
		}).Info("Passing through middleware")
		apiKey := r.Header.Get("x-api-key")

		if apiKey != "" {
			ch := make(chan APIKeyContext)
			log.WithFields(log.Fields{
				"middleware": v.GetID(),
				"key":        apiKey,
			}).Info("Attempting to verify API Key")
			form := url.Values{}
			form.Add("scopeuuid", v.Proxy.ProxyScope)
			form.Add("key", apiKey)
			form.Add("uriPath", v.Proxy.Basepath)
			form.Add("action", "verify")
			go func(form url.Values) {

				apiUrl := "http://localhost:9090/verifiers/apikey"
				resp, err := http.PostForm(apiUrl, form)

				if err != nil {
					ch <- APIKeyContext{Type: "RequestError"}
				}

				var apiKeyContext APIKeyContext
				decoder := json.NewDecoder(resp.Body)
				decoder.UseNumber()

				if err := decoder.Decode(&apiKeyContext); err != nil {
					ch <- APIKeyContext{Type: "DecodeError"}
				} else {
					ch <- apiKeyContext
				}
			}(form)

			for {
				select {
				case apiKeyContext := <-ch:
					if verifyError := apiKeyContext.Type; verifyError == "ErrorResult" {
						log.WithFields(log.Fields{
							"middleware": v.GetID(),
							"key":        apiKey,
						}).Info("Key Not Verified")
						http.Error(w, http.StatusText(401), 401)
					} else if verifyError == "DecodeError" || verifyError == "RequestError" {
						http.Error(w, http.StatusText(500), 500)
					} else {
						status := apiKeyContext.Result.Status
						if status == "APPROVED" {
							log.WithFields(log.Fields{
								"middleware": v.GetID(),
							}).Info("Key Verified")
							next.ServeHTTP(w, r)
						} else {
							log.WithFields(log.Fields{
								"middleware": v.GetID(),
								"key":        apiKey,
							}).Info("Key Not Verified")
							http.Error(w, http.StatusText(401), 401)
						}
					}
					return
				}
			}

		} else {
			log.WithFields(log.Fields{
				"middleware": v.GetID(),
			}).Info("No key present to verify")
			http.Error(w, http.StatusText(401), 401)
		}
	})
}
