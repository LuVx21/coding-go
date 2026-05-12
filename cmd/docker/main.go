package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/fmt_x"
	"github.com/luvx21/coding-go/coding-common/nets_x"
	"github.com/luvx21/coding-go/coding-common/os_x"
	"github.com/luvx21/coding-go/coding-common/times_x"
	"github.com/moby/moby/client"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"golang.org/x/time/rate"
)

const (
	urlTemplate = "https://hub.docker.com/v2/repositories/{{.namespace}}/{{.imageName}}/tags"
	//urlTemplate = "https://hub.docker.com/v2/namespaces/{{.namespace}}/repositories/{{.imageName}}/tags"

	// "https://registry-1.docker.io/v2/library/nginx/tags/list"
	tagsListUrl  = "https://%s/v2/%s/tags/list"
	manifestsUrl = "https://%s/v2/%s/manifests/%s"
	tokenUrl     = "%s?service=%s&scope=%s"
)

var (
	rateLimiter = rate.NewLimiter(1, 1)
	brower      = gorequest.New().Timeout(time.Minute)
	rootCmd     = &cobra.Command{Use: "docker-tags", Long: `docker 相关的一些工具命令`}
	tokenMap    sync.Map
)
var (
	image               string
	pageSize, pageCount int
	proxy               string
)

func initFlag() {
	check := &cobra.Command{Use: "check", Short: "检查本地镜像是否存在更新", Run: checkFunc}
	search := &cobra.Command{Use: "search", Short: "查询指定镜像的tag", Run: searchFunc}
	search.Flags().StringVarP(&image, "image", "i", "", "镜像名,foo/bar:v1,可无namespace(默认library)")
	search.Flags().IntVarP(&pageSize, "ps", "", 25, "单页查询数量")
	search.Flags().IntVarP(&pageCount, "pc", "", 1, "查询多少页")

	sunCmds := []*cobra.Command{check, search}
	for _, cmd := range sunCmds {
		cmd.Flags().StringVarP(&proxy, "proxy", "p", "", "代理地址,如http://127.0.0.1:7890, 会读取http_proxy环境变量")
	}
	rootCmd.AddCommand(sunCmds...)
}

func main() {
	initFlag()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func checkFunc(cmd *cobra.Command, args []string) {
	cli, err := client.New(
		client.WithHost("unix:///var/run/docker.sock"),
		client.WithUserAgent("my-application/1.0.0"),
	)
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	_, err = cli.Ping(context.Background(), client.PingOptions{})
	if err != nil {
		fmt.Printf("Docker daemon 连接失败: %v\n", err)
		os.Exit(1)
	}
	images, _ := cli.ImageList(context.Background(), client.ImageListOptions{All: true, Manifests: true})
	rows := make([]table.Row, 0)
	for _, is := range images.Items {
		iir, _ := cli.ImageInspect(context.Background(), is.ID)
		if len(iir.RepoTags) == 0 || len(iir.RepoDigests) == 0 {
			continue
		}
		temp := parseImage(iir.RepoTags[0])
		temp.Full = iir.RepoTags[0]
		temp.Tag = pickTag(iir.RepoTags[0])
		temp.Digest = pickDigest(iir.RepoDigests[0])
		temp.Os = iir.Os
		temp.Architecture = iir.Architecture
		temp.Size = iir.Size

		fmt_x.Infoln(temp.Full, "......")
		validVersions := []*Version{{Orinal: temp.Tag}}
		if temp.Tag != "latest" {
			baseVersion, formatStr, ok := FromTag(temp.Tag)
			if !ok {
				fmt_x.Warningln("Tag解析失败", temp.Full)
				continue
			}
			remoteTags := getLatestTags(setProxy(brower), temp)
			for _, rTag := range remoteTags {
				rt := rTag.String()
				version, tagFormatStr, ok := FromTag(rt)
				if !ok || !formatCompatible(baseVersion, version, formatStr, tagFormatStr) {
					continue
				}
				validVersions = append(validVersions, version)
			}
			validVersions = deduplicateVersions(validVersions)
			sort.Slice(validVersions, func(i, j int) bool {
				return validVersions[i].Compare(validVersions[j]) < 0
			})
		}

		temp.NewerTag = validVersions[len(validVersions)-1].Orinal
		if temp.NewerTag == temp.Tag {
			temp.NewerTag = ""
			digest := getLatestDigest(setProxy(brower), temp)
			if (digest != "") && (digest != temp.Digest) {
				temp.NewerTag = temp.Tag
			}
		}
		rows = append(rows, table.Row{temp.Full, common_x.IfThen(temp.NewerTag == "", temp.Tag, temp.Tag+" -> "+temp.NewerTag)})
	}
	print(table.Row{"Reference", "Status"}, rows)
}

func getLatestTags(brower *gorequest.SuperAgent, imageInfo ImageInfo) []gjson.Result {
	token, exist := tokenMap.Load(imageInfo.Full)
	var (
		r       gorequest.Response
		data    []gjson.Result
		getTags = func(token string) (gorequest.Response, []gjson.Result) {
			rateLimiter.Wait(context.TODO())
			c := brower.Get(fmt.Sprintf(tagsListUrl, imageInfo.Registry, imageInfo.Namespace+"/"+imageInfo.Image))
			if token != "" {
				c = c.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			}
			a, b, errs := c.SetCurlCommand(false).End()
			if len(errs) == 0 && a.StatusCode == http.StatusOK {
				return a, gjson.Get(b, "tags").Array()
			}
			return a, nil
		}
	)
	r, data = getTags(common_x.IfThen(exist, cast_x.ToString(token), ""))
	if r == nil {
		return nil
	}
	if r.StatusCode == http.StatusOK {
		return data
	}
	wwwHeader := r.Header["Www-Authenticate"]
	if len(wwwHeader) == 0 {
		return nil
	}
	realm, service, scope := parseWwwAuthenticateManual(wwwHeader[0])
	rateLimiter.Wait(context.TODO())
	r, b, errs := brower.Get(fmt.Sprintf(tokenUrl, realm, service, scope)).SetCurlCommand(false).End()
	if len(errs) > 0 {
		return nil
	}
	token = gjson.Get(b, "token").String()
	tokenMap.Store(imageInfo.Full, token)
	r, data = getTags(token.(string))
	return common_x.IfThen(r.StatusCode == http.StatusOK, data, nil)
}

func getLatestDigest(brower *gorequest.SuperAgent, imageInfo ImageInfo) string {
	token, exist := tokenMap.Load(imageInfo.Full)
	var (
		r       gorequest.Response
		data    string
		getTags = func(token string) (gorequest.Response, string) {
			rateLimiter.Wait(context.TODO())
			c := brower.Get(fmt.Sprintf(manifestsUrl, imageInfo.Registry, imageInfo.Namespace+"/"+imageInfo.Image, imageInfo.Tag))
			if token != "" {
				c = c.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			}
			a, _, errs := c.SetCurlCommand(false).End()
			if len(errs) == 0 && a.StatusCode == http.StatusOK {
				return a, a.Header["Docker-Content-Digest"][0]
			}
			return a, ""
		}
	)
	r, data = getTags(common_x.IfThen(exist, cast_x.ToString(token), ""))
	if r == nil {
		return ""
	}
	if r.StatusCode == http.StatusOK {
		return data
	}
	wwwHeader := r.Header["Www-Authenticate"]
	if len(wwwHeader) == 0 {
		return ""
	}
	realm, service, scope := parseWwwAuthenticateManual(wwwHeader[0])
	rateLimiter.Wait(context.TODO())
	r, b, errs := brower.Get(fmt.Sprintf(tokenUrl, realm, service, scope)).SetCurlCommand(false).End()
	if len(errs) > 0 {
		return ""
	}
	token = gjson.Get(b, "token").String()
	tokenMap.Store(imageInfo.Full, token)
	r, data = getTags(token.(string))
	return common_x.IfThen(r.StatusCode == http.StatusOK, data, "")
}
func searchFunc(cmd *cobra.Command, args []string) {
	if len(image) == 0 {
		fmt_x.Warningln("必须指定镜像名称")
		os.Exit(0)
	}

	info := parseImage(image)
	namespace, name, tag := info.Namespace, info.Image, common_x.IfThen(info.Tag == default_tag, "", info.Tag)
	paths := map[string]any{
		"namespace": namespace,
		"imageName": name,
	}
	urlStr, _ := nets_x.UrlAddPath(urlTemplate, paths)

	rows := make([]table.Row, 0)
	for page := 1; page <= pageCount; page++ {
		m := map[string]any{
			"page":      page,
			"page_size": pageSize,
			"ordering":  "last_updated",
			"name":      tag,
		}
		pUrl, _ := nets_x.UrlAddQuery(urlStr, m)
		_ = rateLimiter.Wait(context.TODO())
		fmt.Println("请求地址:", pUrl.String())
		_, body, _ := setProxy(brower).Get(pUrl.String()).End()
		for _, result := range gjson.Get(body, "results").Array() {
			for _, _image := range result.Get("images").Array() {
				os := _image.Get("os").Str
				architecture := _image.Get("architecture").Str
				if os == "unknown" || architecture == "unknown" {
					continue
				}
				date, _ := times_x.StringToDate(_image.Get("last_pushed").Str)
				rows = append(rows, table.Row{
					result.Get("name").Str,
					os + "/" + architecture,
					date.Format(time.DateTime),
					strconv.FormatUint(cast_x.ToUint64(_image.Get("size").Raw)/1024/1024, 10) + " MB",
					_image.Get("digest").Str,
				})
			}
		}
	}
	print(table.Row{"Tag", "Arch", "Pushed", "Size", "Digest"}, rows)
}

func setProxy(*gorequest.SuperAgent) *gorequest.SuperAgent {
	if len(proxy) == 0 {
		proxy = os_x.Getenv("http_proxy")
	}
	if len(proxy) > 0 {
		brower = brower.Proxy(proxy)
	}
	return brower
}

func print(header table.Row, rows []table.Row) {
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	t := table.NewWriter()
	t.SetTitle(image)
	t.AppendHeader(header)
	for _, row := range rows {
		t.AppendRow(row, rowConfigAutoMerge)
	}
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true, Align: text.AlignLeft},
	})
	//t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateRows = true
	fmt.Println(t.Render())
}
