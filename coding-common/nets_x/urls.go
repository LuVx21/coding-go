package nets_x

import (
    "github.com/luvx21/coding-go/coding-common/cast_x"
    "net/url"
)

func UrlAddQuery(urlStr string, queryMap map[string]any) (*url.URL, error) {
    pUrl, _ := url.Parse(urlStr)
    query := pUrl.Query()
    for k, v := range queryMap {
        query.Set(k, cast_x.ToString(v))
    }
    pUrl.RawQuery = query.Encode()
    return pUrl, nil
}
