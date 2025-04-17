package soup

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/luvx21/coding-go/coding-common/retry"
	"github.com/luvx21/coding-go/coding-common/slices_x"
	"github.com/luvx21/coding-go/coding-common/strings_x"
	"github.com/luvx21/coding-go/infra/logs"
	"luvx/gin/common/consts"
)

func (param SpiderParam) content(title, _url string) PageContent {
	paramConfig := param.Content
	contentList := make([]string, 0, 16)
	var pubDate string
	categorySet := make([]string, 0, 2)
	pageUrl := _url
	for i := 0; !strings_x.IsBlank(pageUrl); i++ {
		unescape, _ := url.QueryUnescape(pageUrl)
		logs.Log.Infof("解析内容页: No.%d %s-%s\n", i+1, title, unescape)
		doc, err := request("解析内容页重试", pageUrl)
		if err != nil {
			logs.Log.Error(err)
			break
		}
		if i == 0 {
			contentTitleRule := paramConfig.ContentTitleRule
			if contentTitleRule.Valid() {
				title = getValue1(doc.Selection, contentTitleRule)
			}
			rule := paramConfig.ContentPubDateRule
			if rule.Valid() {
				pubDate = getValue1(doc.Selection, rule)
			}
			contentCategoryRuleList := paramConfig.ContentCategoryRuleList
			if empty, rules := slices_x.IsEmpty(contentCategoryRuleList); !empty {
				for _, _rule := range rules {
					afterSelect := getValueListAfterSelect(doc.Selection, _rule)
					categorySet = append(categorySet, afterSelect...)
				}
			}
		}

		rule := paramConfig.ContentRule
		doc.Find(rule.ElementQuery).Each(func(i int, s *goquery.Selection) {
			value := getValue(s, rule.ValueQuery)
			contentList = append(contentList, value)
		})

		pageUrl = ""
		urlRule := paramConfig.ContentNextPageUrlRule
		if urlRule.Valid() {
			value := getValue1(doc.Selection, urlRule)
			if len(value) > 0 {
				domain := urlAddDomain(_url, value)
				uu := paramConfig.ContentNextPageUrlPostProcessor(domain)
				if pageUrl != uu {
					pageUrl = uu
				}
			}
		}
	}
	page := PageContent{
		Id:          consts.IdWorker.NextId(),
		SpiderKey:   param.SpiderKey,
		Url:         _url,
		Title:       title,
		PubDate:     pubDate,
		CategorySet: categorySet,
		Content:     contentList,
		Invalid:     0,
		CreateTime:  time.Now(),
	}
	return paramConfig.ContentPostProcessor(page)
}

func (param SpiderParam) Visit() []PageContent {
	paramConfig := param.Index
	result := make([]PageContent, 0, 10)

	_url := param.StartUrl
	pageUrl := _url
	for i := 0; i < paramConfig.PageCount && len(pageUrl) > 0; i++ {
		logs.Log.Infoln("解析目录页:", pageUrl)
		doc, err := request("解析目录页重试", pageUrl)
		if err != nil {
			logs.Log.Error(err)
			break
		}
		temp1 := make([]*goquery.Selection, 0, 10)
		doc.Find(paramConfig.IndexItemListRule).Each(func(i int, s *goquery.Selection) {
			temp1 = append(temp1, s)
		})

		indexList := paramConfig.IndexItemListPostProcessor(temp1)
		if len(indexList) == 0 {
			return result
		}
		_max := min(len(indexList), paramConfig.CountInPage)
		for k := range _max {
			element := indexList[k]
			title := strings.TrimSpace(getValue1(element, paramConfig.IndexItemTitleRule))
			href := getValue1(element, paramConfig.IndexItemUrlRule)
			href = urlAddDomain(_url, href)
			logs.Log.Debugf("目录页内容:No.%d %s %s\n", k+1, title, href)
			if strings_x.IsBlank(href) || slices.Contains(paramConfig.IgnoreIndexItemUrlSet, href) {
				continue
			}
			content := param.content(title, href)
			logs.Log.Debug("内容页内容:\n", content.Content)
			if len(content.Content) == 0 {
				logs.Log.Warnln("内容页内容为空")
				continue
			}
			paramConfig.IgnoreIndexItemUrlSet = append(paramConfig.IgnoreIndexItemUrlSet, content.Url)
			result = append(result, content)
		}

		pageUrl = ""
		rule := paramConfig.IndexNextPageUrlRule
		if rule.Valid() {
			value := getValue1(doc.Selection, rule)
			if len(value) > 0 {
				domain := urlAddDomain(_url, value)
				uu := paramConfig.IndexNextPageUrlPostProcessor(domain)
				if pageUrl != uu {
					pageUrl = uu
				}
			}
		}
	}

	return result
}

func request(name, pageUrl string) (*goquery.Document, error) {
	f := func() *goquery.Document {
		limiter := consts.GetRateLimiter(pageUrl)
		_ = limiter.Wait(context.Background())
		dd, e := goquery.NewDocument(pageUrl)
		if e != nil {
			panic(e)
		}
		return dd
	}
	return retry.SupplyWithRetry(name, f, 5, 5*time.Second)
}

func getValueListAfterSelect(element *goquery.Selection, rule QueryRule) []string {
	elementQuery := rule.ElementQuery
	split := strings.SplitSeq(elementQuery, "|")
	for eq := range split {
		result := make([]string, 0, 8)
		element.Find(eq).Each(func(i int, s *goquery.Selection) {
			value := getValue(s, rule.ValueQuery)
			result = append(result, value)
		})
		if len(result) > 0 {
			return result
		}
	}
	return nil
}

func getValue1(element *goquery.Selection, rule QueryRule) string {
	if !rule.Valid() {
		return ""
	}
	elementQuery := rule.ElementQuery
	if strings_x.IsBlank(elementQuery) {
		return getValue(element, rule.ValueQuery)
	}
	split := strings.SplitSeq(elementQuery, "|")
	for eq := range split {
		first := element.Find(eq).First()
		if first == nil {
			continue
		}
		value := getValue(first, rule.ValueQuery)
		if !strings_x.IsEmpty(value) {
			return value
		}
	}
	return ""
}

func getValue(element *goquery.Selection, valueQuery string) string {
	if strings_x.IsBlank(valueQuery) {
		return ""
	}

	split := strings.Split(valueQuery, "|")
	var result string
	for _, q := range split {
		switch q {
		case "text":
			result = element.Text()
		case "data", "href":
			result, _ = element.Attr(q)
		default:
			result, _ = element.Attr(q)
		}
		if !strings_x.IsEmpty(result) {
			break
		}
	}
	return result
}

func urlAddDomain(baseUrl, urlWithoutDomain string) string {
	if !strings.HasPrefix(urlWithoutDomain, "http") {
		u, _ := url.Parse(baseUrl)
		urlWithoutDomain = fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, urlWithoutDomain)
	}
	return urlWithoutDomain
}
