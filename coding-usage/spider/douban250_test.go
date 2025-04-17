package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"testing"
	"time"
)

func Test_01(t *testing.T) {
	url := ""
	//url := ""
	doc, _ := goquery.NewDocument(url)

	doc.Find("head meta[property=\"article:published_time\"]").Each(func(i int, s *goquery.Selection) {
		val, exists := s.Attr("content")
		if exists {
			fmt.Println(val)
		}
	})

	doc.Find("nav.elementor-pagination a.prev").Each(func(i int, s *goquery.Selection) {
		val, exists := s.Attr("href")
		fmt.Println(val, exists, s.Text())
	})
}

func Test02(t *testing.T) {
	start := time.Now()
	ch := make(chan bool)
	for i := range 10 {
		go parseUrls("https://movie.douban.com/top250?start="+strconv.Itoa(25*i), ch)
	}

	for range 10 {
		<-ch
	}

	elapsed := time.Since(start)
	fmt.Printf("Took %s", elapsed)
}
