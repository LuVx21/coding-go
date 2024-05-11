package libsql

import "fmt"

const (
    Db    = "main-luvx21"
    Token = ""
    url   = "libsql://%s.turso.io"
)

var (
    Url = fmt.Sprintf(url, Db)
)
