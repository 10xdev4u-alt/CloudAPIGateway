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

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	
	// Customize the director to log upstream requests if needed
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		logger.Log.Debug("Forwarding request",
			zap.String("target_host", req.URL.Host),
			zap.String("target_path", req.URL.Path),
		)
	}

	return proxy, nil
}
