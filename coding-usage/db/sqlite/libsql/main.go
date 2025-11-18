package libsql

import (
	"fmt"
	"github.com/luvx21/coding-go/coding-usage/db"
)

var (
	Db    = "main-luvx21"
	Token = db.Token
	url   = "libsql://%s.turso.io"
)

var (
	Url = fmt.Sprintf(url, Db)
)
