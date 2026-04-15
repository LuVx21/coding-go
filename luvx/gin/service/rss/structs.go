package rss

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/luvx21/coding-go/coding-common/slices_x"
)

const (
	XmlFormat = `<?xml version="1.0" encoding="UTF-8"?>
<rss xmlns:atom="http://www.w3.org/2005/Atom" version="2.0">
    <channel>
        <title><![CDATA[%s]]></title>
%s
    </channel>
</rss>
`
)

type RSS struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		PubDate     string    `xml:"pubDate"`
		Items       []RssItem `xml:"item"`
	} `xml:"channel"`
}

func (aa *RSS) ToString() string {
	feeds := slices_x.Transfer(func(item RssItem) *RssItem { return &item }, aa.Channel.Items...)
	return ToRssXml(feeds, aa.Channel.Title)
}

type RssItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
	Link        string   `xml:"link"`
	Guid        string   `xml:"guid"`
	Author      string   `xml:"author"`
	Categories  []string `xml:"category,omitempty"`
}

func (aa *RssItem) ToString() string {
	c := ""
	if len(aa.Categories) > 0 {
		c = strings.Join(slices_x.Transfer(func(s string) string { return "    <category><![CDATA[ " + s + " ]]></category>" }, aa.Categories...), "\n")
	}

	return fmt.Sprintf(`
<item>
    <title>
        <![CDATA[ %v ]]>
    </title>
    <description>
        <![CDATA[ %v ]]>
    </description>
    <pubDate>%v</pubDate>
    <link>%v</link>
    <guid>%v</guid>
    <author>
        <![CDATA[ %v ]]>
    </author>
%s
</item>`, aa.Title, aa.Description, aa.PubDate, aa.Link, aa.Guid, aa.Author, c)
}

func ToRssXml(items []*RssItem, title string) string {
	if title == "" {
		title = "RssFeed"
	}
	ss := []string{}
	if len(items) > 0 {
		ss = slices_x.Transfer(func(item *RssItem) string {
			return item.ToString()
		}, items...)
		// ss = parallel.Map(items, func(item *RssItem, index int) string {
		// 	return item.ToString()
		// })
	}
	return fmt.Sprintf(XmlFormat, title, strings.Join(ss, "\n"))
}
