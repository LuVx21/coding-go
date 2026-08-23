package word

import (
	"fmt"
	"testing"

	"github.com/fumiama/jieba"
	"github.com/go-ego/gse"
	"github.com/go-ego/gse/hmm/extracker"
	"github.com/luvx21/coding-go/coding-common/fmt_x"
)

var (
	words = []string{
		"我来到北京清华大学",
		"小明硕士毕业于中国科学院计算所，后在日本京都大学深造",
		"Go语言高性能分词库gse非常好用",
		"大模型推理优化是自然语言处理工程落地中的重要研究方向，向量检索与语义召回显著提升搜索效果",
	}
)

func Test_gse_01(t *testing.T) {
	var seg gse.Segmenter
	seg.LoadDict() // 加载默认字典
	seg.LoadStop() // 加载内置停用词

	var te extracker.TagExtracter
	te.WithGse(seg)
	te.LoadIdf() // 加载 IDF 词典

	for _, text := range words {
		fmt_x.PrintlnRow0("精确模式", text, seg.Trim(seg.CutStop(text, true)))
		// fmt_x.PrintlnRow0("搜索引擎模式", text, seg.CutSearch(text, true))
		tags := te.ExtractTags(text, 3)
		for _, t := range tags {
			fmt.Printf("%s\t%.4f\n", t.Text, t.Weight)
		}
	}
}

func Test_jieba_01(t *testing.T) {
	seg, _ := jieba.LoadDictionaryAt("dict.txt")
	for _, text := range words {
		fmt_x.PrintlnRow0("精确模式", text, seg.Cut(text, true))
	}
}
