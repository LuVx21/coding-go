package libsql

import (
	"fmt"

	"github.com/luvx21/coding-go/coding-usage/db"
)

var (
	Url = fmt.Sprintf("libsql://%s.turso.io", db.Db)
)
