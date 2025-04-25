package infra_libsql

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_exec_00(t *testing.T) {
	rs := ExecSql("luvx21", "main", true, "select * from user;", "select * from common_key_value;")

	for _, r := range rs {
		bytes, _ := json.MarshalIndent(ParseResult(r.Response.Result), "", "  ")
		fmt.Println(string(bytes))
	}
}

func Test_exec_01(t *testing.T) {
	rs := ExecSqlArg("luvx21", "main", true,
		"insert into user('user_name', 'password', 'age', 'birthday', 'update_time', 'ext') values(?, ?, ?, ?, ?, ?);",
		Arg{Type: "text", Value: "foo"},
		Arg{Type: "text", Value: "bar"},
		Arg{Type: "integer", Value: "18"},
		Arg{Type: "text", Value: "1992-10-21 00:01:02"},
		Arg{Type: "text", Value: "2025-01-01 00:01:02"},
		Arg{Type: "text", Value: "{\"b\": 2}"},
	)

	for _, r := range rs {
		if r.Type == "ok" {
			fmt.Println(r.Response)
		} else if r.Type == "error" {
			fmt.Println(r.Error)
		}
	}
}
