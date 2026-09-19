package docker

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/luvx21/coding-go/coding-common/nets_x"
	"github.com/parnurzeal/gorequest"
	"github.com/tidwall/gjson"
)

const (
	socket      = "/var/run/docker.sock"
	api_version = "1.54"
)

var (
	a = "http://127.0.0.1/v" + api_version + "/containers/json"
	b = "http://127.0.0.1/v" + api_version + "/containers/{{.ContainerId}}/json"

	request = client()
)

func Containers(all bool) {
	pUrl, _ := nets_x.UrlAddQuery(a, map[string]any{"all": all})
	_, body, _ := request.Get(pUrl.String()).
		Set("Content-Type", "application/json").
		EndBytes()
	for _, c := range gjson.ParseBytes(body).Array() {
		fmt.Println(
			c.Get("Names").Array()[0].String(),
			c.Get("NetworkSettings.Networks.net_common.IPAddress").String(),
		)
	}
}

func Container(id string) {
	nets_x.UrlAddPath(b, map[string]any{"ContainerId": id})
	_, body, _ := request.Get(b).
		Set("Content-Type", "application/json").
		EndBytes()
	_ = gjson.ParseBytes(body)
}

func client() *gorequest.SuperAgent {
	customTransport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := net.Dialer{Timeout: 5 * time.Second}
			return dialer.DialContext(ctx, "unix", socket)
		},
	}
	request := gorequest.New()
	request.Transport = customTransport
	return request
}
