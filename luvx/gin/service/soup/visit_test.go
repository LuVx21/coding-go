package soup

import (
    "fmt"
    "testing"
)

func Test_00(t *testing.T) {
    _json := `
{
    "index": {
        "pageCount": 1,
        "countInPage": 1,
        "indexItemUrlRule": {
            "valueQuery": "href",
            "elementQuery": "a"
        },
        "indexItemListRule": "div.article div.card",
        "indexItemTitleRule": {
            "valueQuery": "text",
            "elementQuery": "a span.card-title"
        },
        "indexNextPageUrlRule": {
            "valueQuery": "href",
            "elementQuery": "div.paging a.left"
        },
        "indexItemListPostProcessor": "reversed"
    },
    "content": {
      "contentRule": {
          "valueQuery": "text",
          "elementQuery": "div.article-card-content div#articleContent"
      }
    },
    "startUrl": "https://xxx.me/page/14/",
    "spiderKey": "xxx"
}
`
    of := Of(_json)
    fmt.Printf("%#v\n", of)
    visit := of.Visit()
    fmt.Println(visit[0].Content)
}
