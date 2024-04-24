package libsql

import "fmt"

const (
    db    = ""
    Token = ""
    url   = "libsql://%s.turso.io"
)

var (
    Url = fmt.Sprintf(url, db)
)
