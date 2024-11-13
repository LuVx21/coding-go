package main

import (
    "context"
    "flag"
    "fmt"
    "github.com/jedib0t/go-pretty/v6/table"
    "github.com/jedib0t/go-pretty/v6/text"
    "github.com/luvx21/coding-go/coding-common/nets_x"
    "github.com/parnurzeal/gorequest"
    "github.com/tidwall/gjson"
    "golang.org/x/time/rate"
    "strings"
    "time"
)

var (
    image     = flag.String("image", "", "镜像名,foo/bar:v1,可无namespace(默认library)")
    pageSize  = flag.Int("ps", 25, "单页查询数量")
    pageCount = flag.Int("pc", 1, "查询多少页")
    proxy     = flag.String("proxy", "", "代理地址")

    rateLimiter = rate.NewLimiter(1, 1)
    client      = gorequest.New().Timeout(time.Minute)
)

const (
    urlTemplate = "https://hub.docker.com/v2/repositories/{{.namespace}}/{{.imageName}}/tags"
    //urlTemplate = "https://hub.docker.com/v2/namespaces/{{.namespace}}/repositories/{{.imageName}}/tags"
)

func main() {
    flag.Parse()

    namespace, name, tag := imageName(*image)
    paths := map[string]any{
        "namespace": namespace,
        "imageName": name,
    }
    urlStr, _ := nets_x.UrlAddPath(urlTemplate, paths)

    rows := make([]table.Row, 0)
    for page := 1; page <= *pageCount; page++ {
        m := map[string]any{
            "page":      page,
            "page_size": pageSize,
            "ordering":  "last_updated",
            "name":      tag,
        }
        pUrl, _ := nets_x.UrlAddQuery(urlStr, m)
        if len(*proxy) != 0 {
            client = client.Proxy(*proxy)
        }
        _ = rateLimiter.Wait(context.TODO())
        fmt.Println("请求地址:", pUrl.String())
        _, body, _ := client.Get(pUrl.String()).End()
        for _, result := range gjson.Get(body, "results").Array() {
            for _, _image := range result.Get("images").Array() {
                os := _image.Get("os").Str
                architecture := _image.Get("architecture").Str
                if os == "unknown" || architecture == "unknown" {
                    continue
                }
                rows = append(rows, table.Row{
                    result.Get("name").Str,
                    os + "/" + architecture,
                    _image.Get("digest").Str,
                })
            }
        }
    }

    rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

    t := table.NewWriter()
    t.SetTitle(*image)
    t.AppendHeader(table.Row{"Tag", "Arch", "Digest"})
    for _, row := range rows {
        t.AppendRow(row, rowConfigAutoMerge)
    }
    t.SetColumnConfigs([]table.ColumnConfig{
        {Number: 1, AutoMerge: true, Align: text.AlignRight},
    })
    //t.SetStyle(table.StyleLight)
    t.Style().Options.SeparateRows = true
    fmt.Println(t.Render())
}

func imageName(image string) (namespace, name, tag string) {
    arr1 := strings.Split(image, ":")

    if len(arr1) <= 1 {
        tag = ""
    } else {
        tag = arr1[1]
    }
    arr2 := strings.Split(arr1[0], "/")
    if len(arr2) == 1 {
        namespace = "library"
        name = arr2[0]
    } else {
        namespace = arr2[0]
        name = arr2[1]
    }
    return
}
