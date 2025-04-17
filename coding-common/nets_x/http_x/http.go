package http_x

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/luvx21/coding-go/coding-common/os_x"
)

const (
	User_Agent = "User-Agent"
	Accept     = "Accept"
)

func HttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	proxyUrl := os_x.Getenv(os_x.ENV_HTTP_PROXY)
	if proxyUrl != "" {
		_proxyUrl, _ := url.Parse(proxyUrl)
		tr.Proxy = http.ProxyURL(_proxyUrl)
	}

	jar, _ := cookiejar.New(nil)
	return &http.Client{
		Transport: tr,
		Jar:       jar,
		Timeout:   30 * time.Second,
	}
}

func SetHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}
