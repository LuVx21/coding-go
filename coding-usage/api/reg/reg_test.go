package reg

import (
	"github.com/luvx21/coding-go/coding-common/fmt_x"
	"github.com/luvx21/coding-go/coding-common/slices_x"
	"regexp"
	"testing"
)

func Test_00(t *testing.T) {
	s := "<a href=\"https://localhost:8080/0.jpg\">查看图片</a> 测试测试 <a href=\"https://localhost:8080/1.jpg\">查看图片</a>"

	sampleRegexp := regexp.MustCompile(`<a\s+[^>]*href="(.*?)".*?>(.*?)<\/a>`)
	allString := sampleRegexp.FindAllStringSubmatch(s, -1)

	for _, ss := range allString {
		fmt_x.PrintlnRow(slices_x.ToAnySlice(ss)...)
	}
}
