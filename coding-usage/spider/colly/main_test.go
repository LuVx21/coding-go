package colly

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gocolly/colly"
	"github.com/luvx21/coding-go/infra/logs"
)

func Test01(tt *testing.T) {
	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	c := colly.NewCollector()
	c.WithTransport(t)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36 Edg/129.0.0.0'")
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error:", err)
	})

	c.OnHTML("img", func(e *colly.HTMLElement) {
		url := e.Attr("src")
		if url != "" {
			fmt.Println("Found image:", url)
			resp, _ := http.Get(url)
			defer resp.Body.Close()

			file, _ := os.Create(strconv.FormatInt(time.Now().UnixNano(), 10) + ".jpg")
			defer file.Close()

			io.Copy(file, resp.Body)
			fmt.Println("Image saved to", file.Name())
		}
	})

	// path := ""
	// c.Visit("file://" + path)
	c.Visit("")
}

/*
*
m.lnsjkc.com
*/
func Test03(t *testing.T) {
	c1 := colly.NewCollector()
	c1.OnHTML("div#novelcontent p", func(e *colly.HTMLElement) {
		logs.Log.Infoln(strings.ReplaceAll(e.Text, "。", "。\n"))
	})
	c1.OnHTML("div#novelcontent ul.novelbutton a#pb_next", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasSuffix(href, "_2.html") {
			c1.Visit("https://m.lnsjkc.com" + href)
		}
	})

	c := colly.NewCollector()
	c.OnHTML("ul.chapters li a", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		fmt.Println(e.Text, url)
		if url != "" {
			c1.Visit(url)
		}
	})

	c.Visit("https://m.lnsjkc.com/46/46869_7/")
}
