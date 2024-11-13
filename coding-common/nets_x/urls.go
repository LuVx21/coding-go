package nets_x

import (
    "bytes"
    "github.com/luvx21/coding-go/coding-common/cast_x"
    "net/url"
    "text/template"
)

func UrlAddPath(urlTemplate string, pathMap map[string]any) (string, error) {
    if len(urlTemplate) == 0 || len(pathMap) == 0 {
        return urlTemplate, nil
    }

    tmpl, err := template.New("url").Parse(urlTemplate)
    if err != nil {
        return urlTemplate, err
    }
    var result bytes.Buffer
    if err = tmpl.Execute(&result, pathMap); err != nil {
        return urlTemplate, err
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
