package infra_libsql

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_exec_00(t *testing.T) {
	rs := ExecSql("luvx21", "main", true, "select * from user;", "select * from common_key_value;")

	for _, r := range rs {
		bytes, _ := json.MarshalIndent(ParseResult(r), "", "  ")
		fmt.Println(string(bytes))
	}
}
