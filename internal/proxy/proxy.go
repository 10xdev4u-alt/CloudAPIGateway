package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/princetheprogrammer/cloud-api-gateway/internal/logger"
	"go.uber.org/zap"
)

func NewReverseProxy(target string) (*httputil.ReverseProxy, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(targetURL)
			pr.Out.Host = targetURL.Host
			logger.Log.Debug("Forwarding request",
				zap.String("target_host", pr.Out.URL.Host),
				zap.String("target_path", pr.Out.URL.Path),
			)
		},
	}

	return proxy, nil
}
