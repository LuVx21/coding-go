package main

import (
	"github.com/luvx21/coding-go/coding-common/fmt_x"
	"github.com/yanyiwu/gojieba"
	"regexp"
	"strings"
	"testing"
)

func Test_jieba_00(t *testing.T) {
	var s string
	var words []string
	use_hmm := true
	x := gojieba.NewJieba()
	defer x.Free()

	s = "我来到北京清华大学"
	words = x.CutAll(s)
	fmt_x.PrintlnRow0("全模式", s, strings.Join(words, "/"))

	words = x.Cut(s, use_hmm)
	fmt_x.PrintlnRow0("精确模式", s, strings.Join(words, "/"))

	s = "比特币"
	words = x.Cut(s, use_hmm)
	fmt_x.PrintlnRow0("精确模式", s, strings.Join(words, "/"))

	x.AddWord("比特币")
	// `AddWordEx` 支持指定词语的权重，作为 `AddWord` 权重太低加词失败的补充。
	// `tag` 参数可以为空字符串，也可以指定词性。
	// x.AddWordEx("比特币", 100000, "")
	s = "比特币"
	words = x.Cut(s, use_hmm)
	fmt_x.PrintlnRow0("添加词典后,精确模式", s, strings.Join(words, "/"))

	s = "他来到了网易杭研大厦"
	words = x.Cut(s, use_hmm)
	fmt_x.PrintlnRow0("新词识别", s, strings.Join(words, "/"))

	s = "小明硕士毕业于中国科学院计算所，后在日本京都大学深造"
	words = x.CutForSearch(s, use_hmm)
	fmt_x.PrintlnRow0("搜索引擎模式", s, strings.Join(words, "/"))

	s = "长春市长春药店"
	words = x.Tag(s)
	fmt_x.PrintlnRow0("词性标注", s, strings.Join(words, ","))

	s = "区块链"
	words = x.Tag(s)
	fmt_x.PrintlnRow0("词性标注", s, strings.Join(words, ","))

	s = "长江大桥"
	words = x.CutForSearch(s, !use_hmm)
	fmt_x.PrintlnRow0("搜索引擎模式", s, strings.Join(words, "/"))

	wordinfos := x.Tokenize(s, gojieba.SearchMode, !use_hmm)
	fmt_x.PrintlnRow0("Tokenize:(搜索引擎模式)", s, wordinfos)

	wordinfos = x.Tokenize(s, gojieba.DefaultMode, !use_hmm)
	fmt_x.PrintlnRow0("Tokenize:(默认模式)", s, wordinfos)

	keywords := x.ExtractWithWeight(s, 5)
	fmt_x.PrintlnRow0("Extract", s, keywords)
}

func preProcess(text string) string {
	//去除HTML标签
	re := regexp.MustCompile(`(?s)<.*?>`)
	text = re.ReplaceAllString(text, "")
	//去除特殊符号
	re = regexp.MustCompile(`[^\p{Han}\w]+`)
	text = re.ReplaceAllString(text, "")
	//去除停用词
	stopwords := []string{"的", "了", "是", "在", "这", "个"}
	for _, stopword := range stopwords {
		re = regexp.MustCompile(`\b` + stopword + `\b`)
		text = re.ReplaceAllString(text, "")
	}
	return text
}

//func extractKeywords(words []string) []string {
//    tfidf := gojieba.NewTfidf()
//    defer tfidf.Free()
//    for _, word := range words {
//        tfidf.AddDocument([]string{word})
//    }
//    keywords := tfidf.Extract(10)
//    return keywords
//}
