package infra_sql

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/luvx21/coding-go/coding-common/maps_x"
)

func ParseJson2Ddl(tableName, _json string) string {
	_json = strings.TrimSpace(_json)
	var r any
	if !strings.HasPrefix(_json, "{") && !strings.HasPrefix(_json, "[") {
		return ""
	}
	if err := json.Unmarshal([]byte(_json), &r); err != nil {
		return ""
	}

	var m map[string][]string
	switch v := r.(type) {
	case map[string]any:
		m = parseMap(v)
	case []any:
		switch first := v[0].(type) {
		case map[string]any:
			m = parseMap(first)
		case []any:
			return ""
		default:
			return ""
		}
	default:
		return ""
	}
	if len(m) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("CREATE TABLE IF NOT EXISTS `" + tableName + "`\n(\n")
	sb.WriteString("    `id` integer primary key autoincrement,\n")
	fields := maps_x.JoinMapper(m, " ", ",\n", func(k string) string { return "    `" + k + "`" }, func(v []string) string { return v[2] })
	sb.WriteString(fields)
	sb.WriteString("\n);")
	return sb.String()
}

func parseMap(jsonObject map[string]any) map[string][]string {
	r := make(map[string][]string)
	for k, v := range jsonObject {
		var types []string
		switch v.(type) {
		case float64:
			types = []string{"float64", "decimal", "real"}
		case bool:
			types = []string{"bool", "tinyint", "integer"}
		case string:
			types = []string{"string", "text", "text"}
		case map[string]any:
			types = []string{"map[string]any", "json", "text"}
		case []any:
			types = []string{"[]any", "json", "text"}
		case nil:
			types = []string{"nil", "text", "text"}
		}
		r[k] = types
	}
	return r
}

func Test_json_00(t *testing.T) {
	_json := `
	{
		"a":"aaa",
		"b": 1,
		"c": true,
		"d":{},
		"e":[],
		"f":null
	}
	`
	ddl := ParseJson2Ddl("test", _json)
	fmt.Println(ddl)
	ddl = ParseJson2Ddl("test1", "["+_json+"]")
	fmt.Println(ddl)
}
