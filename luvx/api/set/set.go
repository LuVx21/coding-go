package set

import (
    "fmt"
    set "github.com/deckarep/golang-set/v2"
)

func a() {
    // 默认创建的线程安全的，如果无需线程安全
    // 可以使用 NewThreadUnsafeSet 创建，使用方法都是一样的。
    s1 := set.NewSet(1, 2, 3, 4)
    fmt.Printf("3: %t 5:%t \n", s1.Contains(3), s1.Contains(5))

    // interface 参数，可以传递任意类型
    s1.Add("foo")
    s1.Remove(3)
    fmt.Printf("foo: %t 3: %t\n", s1.Contains("foo"), s1.Contains(3))

    s2 := set.NewSet(1, 3, 4, 5)

    // 并集
    fmt.Println(s1.Union(s2))
}
