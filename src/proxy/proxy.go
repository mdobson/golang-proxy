package proxy

import (
	"middleware"
	"net/http"
	"net/http/httputil"
	"net/url"
	"roundtrippers"
)

type ProxyData struct {
	targetURL       string `yaml:"target_url"`
	revision        string `yaml:"revision"`
	targetName      string `yaml:"target_name"`
	proxyName       string `yaml:"name"`
	basepath        string `yaml:"base_path"`
	vhost           string `yaml:"vhost"`
	proxyConfigName string
}

type ReverseProxy struct {
	target      *url.URL
	proxy       *httputil.ReverseProxy
	data        *ProxyData
	middlewares *middleware.RequestMiddlewareSequence
}

func New(proxyConfigName string, proxyData map[string]interface{}) *ReverseProxy {
	proxyStruct := ProxyData{
		targetURL:       proxyData["url"].(string),
		revision:        proxyData["revision"].(string),
		targetName:      proxyData["target_name"].(string),
		proxyName:       proxyData["proxy_name"].(string),
		basepath:        proxyData["base_path"].(string),
		vhost:           proxyData["vhost"].(string),
		proxyConfigName: proxyConfigName,
	}
	url, _ := url.Parse(proxyStruct.targetURL)
	p := httputil.NewSingleHostReverseProxy(url)
	p.Transport = &roundtrippers.ResponseInterceptTransport{http.DefaultTransport}
	middlewareSequence := middleware.CreateService()

	return &ReverseProxy{target: url, proxy: p, data: &proxyStruct, middlewares: &middlewareSequence}
}

func (p *ReverseProxy) FinalMiddleware(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}

func (p *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	finalHandler := http.HandlerFunc(p.FinalMiddleware)
	middlewares := []string{"HeaderSet", "BodyRewrite", "TriggerBadRequest"}
	p.middlewares.Compile(middlewares, finalHandler).ServeHTTP(w, r)
}

func (p *ReverseProxy) Basepath() string {
	return p.data.basepath
}
