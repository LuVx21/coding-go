package nets_x

import (
	"bytes"
	"net/url"
	"text/template"

	"github.com/luvx21/coding-go/coding-common/cast_x"
)

func UrlAddPath(ut string, pathMap map[string]any) (string, error) {
	if len(ut) == 0 || len(pathMap) == 0 {
		return ut, nil
	}

	tmpl, err := template.New("url").Parse(ut)
	if err != nil {
		return ut, err
	}
	var result bytes.Buffer
	if err = tmpl.Execute(&result, pathMap); err != nil {
		return ut, err
	}
	return result.String(), nil
}

func UrlAddQuery(urlStr string, queryMap map[string]any) (*url.URL, error) {
	pUrl, _ := url.Parse(urlStr)
	query := pUrl.Query()
	for k, v := range queryMap {
		query.Set(k, cast_x.ToString(v))
	}
	pUrl.RawQuery = query.Encode()
	return pUrl, nil
}
