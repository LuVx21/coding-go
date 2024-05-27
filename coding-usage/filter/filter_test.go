package filter

import (
    "fmt"
    "github.com/bits-and-blooms/bloom/v3"
    "github.com/linvon/cuckoo-filter"
    "testing"
)

func Test_00(t *testing.T) {
    // 初始化一个布谷鸟过滤器
    // 使用半排序桶
    // 	  每个桶包含 4 个指纹, 每个指纹 4 bits
    // 最大存放元素数量 4096
    cf := cuckoo.NewFilter(4, 4, 1<<12, cuckoo.TableTypePacked)

    // 添加一些元素
    cf.Add([]byte(`Hello World`))
    cf.Add([]byte(`Hello Golang`))

    // 检测元素是否存在
    fmt.Printf("%v\n", cf.Contain([]byte(`Hello World`)))
    fmt.Printf("%v\n", cf.Contain([]byte(`Hello Golang`)))
    fmt.Printf("%v\n", cf.Contain([]byte(`Hello Rust`)))

    // 输出元素数量
    fmt.Printf("filter size = %d\n", cf.Size())

    // 删除元素
    cf.Delete([]byte(`Hello World`))

    // 输出过滤器统计信息
    fmt.Println("\n", cf.Info())
}

func Test_01(t *testing.T) {
    // 初始化能够接收 100 万个元素且误判率为 1% 的布隆过滤器
    filter := bloom.NewWithEstimates(1000000, 0.01)

    hw := []byte(`hello world`)
    hg := []byte(`hello golang`)

    filter.Add(hw)

    println(filter.Test(hw)) // true
    println(filter.Test(hg)) // false
}
