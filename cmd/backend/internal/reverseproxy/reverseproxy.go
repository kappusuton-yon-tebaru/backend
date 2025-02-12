package reverseproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
)

type ReverseProxy struct {
	proxy *httputil.ReverseProxy
}

func New(cfg *config.Config) (*ReverseProxy, error) {
	url, err := url.Parse(cfg.Backend.AgentEndpoint)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	return &ReverseProxy{
		proxy,
	}, nil
}

func (rp *ReverseProxy) Forward() gin.HandlerFunc {
	return gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		rp.proxy.ServeHTTP(w, r)
	})
}
