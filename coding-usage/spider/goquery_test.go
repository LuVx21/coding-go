package spider

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/longbridgeapp/opencc"
	"github.com/luvx21/coding-go/coding-usage/common"
	_ "github.com/luvx21/coding-go/infra/logs"
)

var (
	client = &http.Client{}
)

func fetch(url string) *goquery.Document {
	log.Print("发起请求:", url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", common.UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Http get err:", err)
	}
	if resp.StatusCode < 400 {
		log.Print("Http status code:", resp.StatusCode)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func Test_goquery_00(t *testing.T) {
	t2s, _ := opencc.New("t2s")
	file, _ := os.OpenFile(os.ExpandEnv("${HOME}/data/n.txt"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	for _, url := range []string{} {
		doc := fetch(url)

		doc.Find("div#txtcontent0").Each(func(i int, s *goquery.Selection) {
			simplified, _ := t2s.Convert(s.Text())
			fmt.Println(simplified)
			file.WriteString(simplified)
			file.WriteString("\n")
		})
		time.Sleep(time.Second)
	}

	// doc.Find("head meta[property=\"article:published_time\"]").Each(func(i int, s *goquery.Selection) {
	// 	val, exists := s.Attr("content")
	// 	if exists {
	// 		fmt.Println(val)
	// 	}
	// })

	// doc.Find("nav.elementor-pagination a.prev").Each(func(i int, s *goquery.Selection) {
	// 	val, exists := s.Attr("href")
	// 	fmt.Println(val, exists, s.Text())
	// })
}

func Test02(t *testing.T) {
	start := time.Now()
	ch := make(chan bool)
	for i := range 10 {
		go func(url string, ch chan bool) {
			doc := fetch(url)
			doc.Find("ol.grid_view li").Find(".hd").Each(func(index int, ele *goquery.Selection) {
				movieUrl, _ := ele.Find("a").Attr("href")
				fmt.Println(strings.Split(movieUrl, "/")[4], ele.Find(".title").Eq(0).Text())
			})
			time.Sleep(2 * time.Second)
			ch <- true
		}("https://movie.douban.com/top250?start="+strconv.Itoa(25*i), ch)
	}

	for range 10 {
		<-ch
	}

	elapsed := time.Since(start)
	fmt.Printf("Took %s", elapsed)
}
