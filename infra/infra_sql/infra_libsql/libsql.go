package infra_libsql

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/os_x"
	"github.com/luvx21/coding-go/coding-common/slices_x"
	"github.com/parnurzeal/gorequest"
)

const (
	domain       = "%s-%s.turso.io"
	http_dsn     = "https://" + domain
	libsql_dsn   = "libsql://" + domain
	pipeline_url = http_dsn + "/v2/pipeline"

	arg_type_null    = "null"
	arg_type_integer = "integer"
	arg_type_float   = "float"
	arg_type_text    = "text"
	arg_type_blob    = "blob"
)

var (
	CloseRequest  = Request{Type: "close"}
	BeginRequest  = Request{Type: "execute", Stmt: Stmt{Sql: "BEGIN;"}}
	CommitRequest = Request{Type: "execute", Stmt: Stmt{Sql: "COMMIT;"}}

	GoRequest = gorequest.New().Timeout(time.Minute)
)

type Arg struct {
	Type  string `json:"type"` // null, integer, float, text, or blob
	Value string `json:"value"`
}

type Stmt struct {
	Sql       string `json:"sql"`
	Args      []Arg  `json:"args,omitempty"`
	NamedArgs []struct {
		Name  string `json:"name"`
		Value Arg    `json:"value"`
	} `json:"named_args,omitempty"`
}
type Request struct {
	Type string `json:"type"`
	Stmt Stmt   `json:"stmt,omitzero"`
}

func ExecSqlArg(org, dbName string, close bool, sql string, args ...Arg) []ResultMeta {
	request := Request{Type: "execute", Stmt: Stmt{Sql: sql, Args: args}}
	return Exec(org, dbName, close, request)
}

func ExecSql(org, dbName string, close bool, sqls ...string) []ResultMeta {
	stmts := slices_x.Transfer(func(sql string) Request {
		return Request{Type: "execute", Stmt: Stmt{Sql: sql}}
	}, sqls...)
	return Exec(org, dbName, true, stmts...)
}

// Exec close: 本次执行后是否关闭连接
func Exec(org, dbName string, close bool, stmts ...Request) []ResultMeta {
	if !close && len(stmts) == 0 {
		return nil
	}
	new_stmts := stmts
	if close {
		new_stmts = append(new_stmts, CloseRequest)
	}
	_, body, err := GoRequest.Post(fmt.Sprintf(pipeline_url, dbName, org)).
		Set("Authorization", "Bearer "+os_x.Getenv("LIBSQL_TOKEN")).
		Set("Content-Type", "application/json").
		SendMap(map[string][]Request{"requests": new_stmts}).
		End()
	if err != nil {
		slog.Error("发起执行sql异常", "sql", new_stmts)
		return nil
	}

	var r R
	if err := json.Unmarshal([]byte(body), &r); err != nil {
		return nil
	}

	result := make([]ResultMeta, len(stmts))
	for i := range stmts {
		result[i] = r.Results[i]
	}
	return result
}

func ParseResult(r Result) []map[string]any {
	cols, rows := r.Cols, r.Rows
	rr := make([]map[string]any, len(rows))
	for i, row := range rows {
		rowMap := make(map[string]any, len(row))
		for j, cell := range row {
			var v any = cell.Value
			switch cell.Type {
			case "integer":
				v = cast_x.ToInt64(cell.Value)
			case "text":
			case "real":
				v = cast_x.ToFloat64(cell.Value)
			default:
			}
			rowMap[cols[j].Name] = v
		}
		rr[i] = rowMap
	}
	return rr
}

type ColMeta struct {
	Name     string
	Decltype string
}
type Cell struct {
	Type  string // integer text
	Value string
}
type Result struct {
	Cols []ColMeta
	Rows [][]Cell
}
type ResultMeta struct {
	Type     string // ok或error
	Response struct {
		Type   string // execute close
		Result Result
	}
	Error struct { // 执行出错时
		Message string
		Code    string
	}
}
type R struct {
	Results []ResultMeta `json:"results"`
}
