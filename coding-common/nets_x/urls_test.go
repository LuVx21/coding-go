package nets_x

import (
	"fmt"
	"testing"
)

func Test_00(t *testing.T) {
	urlTemplate := "https://hub.docker.com/v2/repositories/{{.namespace}}/{{.imageName}}/tags"
	data := map[string]any{
		"namespace": "luvx",
		"imageName": "jupyter",
	}
	urlStr, _ := UrlAddPath(urlTemplate, data)

	m := map[string]any{
		"page":      1,
		"page_size": 25,
		"ordering":  "last_updated",
		"name":      "latest",
	}
	u, _ := UrlAddQuery(urlStr, m)
	fmt.Println(urlStr, "\n", u.String(), u.Path)
}

func Test_01(t *testing.T) {
}
