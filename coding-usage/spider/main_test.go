package main

import (
    "fmt"
    "github.com/gocolly/colly"
    "github.com/luvx21/coding-go/coding-common/logs"
    "io"
    "net/http"
    "os"
    "strconv"
    "strings"
    "testing"
    "time"
)

func Test01(t *testing.T) {
    c := colly.NewCollector()

    c.OnRequest(func(r *colly.Request) {
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

    c.Visit("")
}

/**
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
