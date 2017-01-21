package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"proxy"
	"strings"

	"github.com/gorilla/mux"

	"archive/zip"

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
	resp, err := http.Get("http://localhost:9000/deployments")

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	json.Unmarshal(body, &latestDeployments)

	for _, d := range latestDeployments {
		fmt.Printf("getting uri: %s\n", d.URI)
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

			obj := make(map[interface{}]interface{})
			err = yaml.Unmarshal([]byte(content), &obj)
			if err != nil {
				panic(err)
			}

			for _, proxyData := range obj {

				m2 := make(map[string]string)
				for key, value := range proxyData.(map[interface{}]interface{}) {
					switch key := key.(type) {
					case string:
						switch value := value.(type) {
						case string:
							m2[key] = value
						}
					}
				}

				p := proxy.New(m2["url"])
				r.Handle(m2["base_path"], http.StripPrefix(m2["base_path"], p))
			}
		}

	}

	log.Fatal(http.ListenAndServe(":8080", r))
}
