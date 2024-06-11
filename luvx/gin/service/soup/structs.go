package soup

import (
    "encoding/json"
    "github.com/PuerkitoBio/goquery"
    "github.com/luvx21/coding-go/coding-common/func_x"
    "github.com/tidwall/gjson"
    "slices"
    "time"
)

type PageContent struct {
    Id          int64  `bson:"_id"`
    SpiderKey   string `bson:"spiderKey"`
    Url         string
    Title       string
    PubDate     string   `bson:"pubDate"`
    CategorySet []string `bson:"categorySet"`
    Content     []string
    Invalid     int
    CreateTime  time.Time `bson:"createTime"`
}

type QueryRule struct {
    ElementQuery string
    ValueQuery   string
}

func (rule QueryRule) Valid() bool {
    return rule.ElementQuery != "" || rule.ValueQuery != ""
}

type index struct {
    PageCount                     int
    CountInPage                   int
    IndexItemListRule             string
    IndexItemListPostProcessor    func([]*goquery.Selection) []*goquery.Selection `json:"-"`
    IndexItemTitleRule            QueryRule
    IndexItemUrlRule              QueryRule
    IgnoreIndexItemUrlSet         []string // Set,去重
    IndexNextPageUrlRule          QueryRule
    IndexNextPageUrlPostProcessor func(string) string `json:"-"`
}
type content struct {
    ContentTitleRule                QueryRule
    ContentPubDateRule              QueryRule
    ContentCategoryRuleList         []QueryRule
    ContentRule                     QueryRule
    ContentPostProcessor            func(PageContent) PageContent `json:"-"`
    ContentNextPageUrlRule          QueryRule
    ContentNextPageUrlPostProcessor func(string) string `json:"-"`
}

type SpiderParam struct {
    SpiderKey string
    StartUrl  string
    Index     index
    Content   content
}

func (param SpiderParam) SetIgnoreIndexItemUrlSet(set []string) SpiderParam {
    param.Index.IgnoreIndexItemUrlSet = set
    return param
}

var (
    indexItemListPostProcessorMap = map[string]func([]*goquery.Selection) []*goquery.Selection{
        "reversed": func(from []*goquery.Selection) []*goquery.Selection {
            slices.Reverse(from)
            return from
        },
    }

    indexNextPageUrlPostProcessorMap   = map[string]func(string) string{}
    contentPostProcessorMap            = map[string]func(PageContent) PageContent{}
    contentNextPageUrlPostProcessorMap = map[string]func(string) string{}
)

func Of(_json string) SpiderParam {
    r := SpiderParam{
        Index: index{
            PageCount:                     2,
            CountInPage:                   2,
            IndexItemListPostProcessor:    func_x.Identity[[]*goquery.Selection],
            IndexNextPageUrlPostProcessor: func_x.Identity[string],
        },
        Content: content{
            ContentPostProcessor:            func_x.Identity[PageContent],
            ContentNextPageUrlPostProcessor: func_x.Identity[string],
        },
    }
    _ = json.Unmarshal([]byte(_json), &r)

    a := gjson.Get(_json, "index.indexItemListPostProcessor").String()
    a2 := gjson.Get(_json, "index.indexNextPageUrlPostProcessor").String()
    a3 := gjson.Get(_json, "content.contentPostProcessor").String()
    a4 := gjson.Get(_json, "content.contentNextPageUrlPostProcessor").String()

    if len(a) > 0 {
        f := indexItemListPostProcessorMap[a]
        if f == nil {
            panic("参数配置错误:indexItemListPostProcessor")
        }
        r.Index.IndexItemListPostProcessor = f
    }
    if len(a2) > 0 {
        f := indexNextPageUrlPostProcessorMap[a2]
        if f == nil {
            panic("参数配置错误:indexNextPageUrlPostProcessor")
        }
        r.Index.IndexNextPageUrlPostProcessor = f
    }
    if len(a3) > 0 {
        f := contentPostProcessorMap[a3]
        if f == nil {
            panic("参数配置错误:contentPostProcessor")
        }
        r.Content.ContentPostProcessor = f
    }
    if len(a4) > 0 {
        f := contentNextPageUrlPostProcessorMap[a4]
        if f == nil {
            panic("参数配置错误:contentNextPageUrlPostProcessor")
        }
        r.Content.ContentNextPageUrlPostProcessor = f
    }
    return r
}
