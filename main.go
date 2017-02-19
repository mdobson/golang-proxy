package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"proxy"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	yaml "gopkg.in/yaml.v2"
)

type deployment struct {
	ID               string          `json:"id"`
	ScopeID          string          `json:"scopeId"`
	Created          string          `json:"created"`
	CreatedBy        string          `json:"createdBy"`
	Updated          string          `json:"updated"`
	UpdatedBy        string          `json:"updatedBy"`
	ConfigJSON       json.RawMessage `json:"configuration"`
	BundleConfigJSON json.RawMessage `json:"bundleConfiguration"`
	DisplayName      string          `json:"displayName"`
	URI              string          `json:"uri"`
}

func main() {
	r := mux.NewRouter()
	type deployments []deployment

	latestDeployments := deployments{}
	apidUrl := "http://localhost:9090/deployments"
	log.WithFields(log.Fields{
		"url": apidUrl,
	}).Info("Downloading deployments from apid")
	resp, err := http.Get(apidUrl)

	if err != nil {
		log.Errorf("Error retrieving deployments: %s", err.Error())
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Errorf("Error reading deployments: %s", err.Error())
		return
	}

	json.Unmarshal(body, &latestDeployments)

	for _, d := range latestDeployments {
		log.Info("Unzipping deployment at URI", d.URI)
		replacedPath := strings.Replace(d.URI, "file://", "", 1)
		zippedFile, err := zip.OpenReader(replacedPath)

		if err != nil {
			panic(err)
		}

		defer zippedFile.Close()

		for _, f := range zippedFile.File {
			rc, err := f.Open()
			if err != nil {
				panic(err)
			}

			content, err := ioutil.ReadAll(rc)

			if err != nil {
				panic(err)
			}

			obj := make(map[string]map[string]interface{})
			err = yaml.Unmarshal([]byte(content), &obj)
			if err != nil {
				panic(err)
			}

			for proxyConfigName, proxyData := range obj {
				p := proxy.New(proxyConfigName, d.ScopeID, proxyData)
				log.WithFields(log.Fields{
					"step":     "deployment",
					"basepath": p.Basepath(),
					"target":   p.Target(),
				}).Info("Proxy found and being deployed")
				r.Handle(fmt.Sprintf("%s{rest:.*}", p.Basepath()), http.StripPrefix(p.Basepath(), p))
			}
		}

	}

	log.Info("Serving http requests at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
