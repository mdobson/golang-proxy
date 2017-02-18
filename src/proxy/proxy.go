package proxy

import (
	"data"
	"middleware"
	"net/http"
	"net/http/httputil"
	"net/url"
	"roundtrippers"
)

type ReverseProxy struct {
	target      *url.URL
	proxy       *httputil.ReverseProxy
	data        *data.ProxyData
	middlewares *middleware.RequestMiddlewareSequence
}

func New(proxyConfigName string, proxyScope string, proxyData map[string]interface{}) *ReverseProxy {
	proxyStruct := data.ProxyData{
		TargetURL:       proxyData["url"].(string),
		Revision:        proxyData["revision"].(string),
		TargetName:      proxyData["target_name"].(string),
		ProxyName:       proxyData["proxy_name"].(string),
		Basepath:        proxyData["base_path"].(string),
		Vhost:           proxyData["vhost"].(string),
		ProxyConfigName: proxyConfigName,
		ProxyScope:      proxyScope,
	}
	url, _ := url.Parse(proxyStruct.TargetURL)
	p := httputil.NewSingleHostReverseProxy(url)
	p.Transport = &roundtrippers.ResponseInterceptTransport{http.DefaultTransport}
	middlewareSequence := middleware.CreateService(proxyStruct)

	return &ReverseProxy{target: url, proxy: p, data: &proxyStruct, middlewares: &middlewareSequence}
}

func (p *ReverseProxy) FinalMiddleware(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}

func (p *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	finalHandler := http.HandlerFunc(p.FinalMiddleware)
	middlewares := []string{"VerifyApiKey", "HeaderSet", "BodyRewrite", "TriggerBadRequest"}
	p.middlewares.Compile(middlewares, finalHandler).ServeHTTP(w, r)
}

func (p *ReverseProxy) Basepath() string {
	return p.data.Basepath
}
