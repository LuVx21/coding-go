package word

import (
	"testing"

	"github.com/fumiama/jieba"
	"github.com/go-ego/gse"
	"github.com/luvx21/coding-go/coding-common/fmt_x"
)

var (
	words = []string{
		"我来到北京清华大学",
		"小明硕士毕业于中国科学院计算所，后在日本京都大学深造",
		"Go语言高性能分词库gse非常好用",
	}
)

func Test_gse_01(t *testing.T) {
	var seg gse.Segmenter
	seg.LoadDict() // 加载默认字典

	for _, text := range words {
		fmt_x.PrintlnRow0("精确模式", text, seg.Cut(text, true))
		fmt_x.PrintlnRow0("搜索引擎模式", text, seg.CutSearch(text, true))
	}
}

func Test_jieba_01(t *testing.T) {
	seg, _ := jieba.LoadDictionaryAt("dict.txt")
	for _, text := range words {
		fmt_x.PrintlnRow0("精确模式", text, seg.Cut(text, true))
	}
}
