package rss

import "fmt"

type ItemRss struct {
    Title       string
    Description string
    PubDate     string
    Link        string
    Guid        int64
    Author      string
}

func (aa *ItemRss) ToString() string {
    return fmt.Sprintf(
        `
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
           </item>
`, aa.Title, aa.Description, aa.PubDate, aa.Link, aa.Guid, aa.Author)
}

func ToRssXml(items []*ItemRss, title string) string {
    if title == "" {
        title = "网络傻事"
    }
    result := ""
    for _, item := range items {
        result += item.ToString()
    }
    s := `<?xml version="1.0" encoding="UTF-8"?>
<rss xmlns:atom="http://www.w3.org/2005/Atom" version="2.0">
    <channel>
        <title><![CDATA[%s]]></title>
%s
    </channel>
</rss>
`
    return fmt.Sprintf(s, title, result)
}
