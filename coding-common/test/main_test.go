package test

import (
	"fmt"
	"testing"
)

var aa = func(name string) func() {
	return BeforeAfterTest(name, func() { fmt.Println("准备阶段") }, func() { fmt.Println("收尾阶段") })
}

func Test_00(t *testing.T) {
	defer aa("Test_00")()

	fmt.Println("进行测试")
}

func Test_01(t *testing.T) {
	defer aa("Test_01")()

	fmt.Println("进行测试")
}
