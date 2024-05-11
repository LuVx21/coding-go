package main

import (
    "bytes"
    "io"
    "os"
    "sync"
    "testing"
    "time"
)

var bufPool = sync.Pool{
    New: func() any {
        return new(bytes.Buffer)
    },
}

func timeNow() time.Time {
    return time.Unix(1136214245, 0)
}

func test(w io.Writer, key, val string) {
    b := bufPool.Get().(*bytes.Buffer) // 类型转换!!!!
    b.Reset()
    b.WriteString(timeNow().UTC().Format(time.RFC3339))
    b.WriteByte(' ')
    b.WriteString(key)
    b.WriteByte('=')
    b.WriteString(val)
    w.Write(b.Bytes())
    bufPool.Put(b)
}
func Test_pool_00(t *testing.T) {
    test(os.Stdout, "path", "/search?q=flowers\n")
}
