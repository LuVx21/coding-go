package jsons

import (
	// json "github.com/bytedance/sonic"
	"encoding/json"
)

func JsonStringToObject(s string, v any) error {
	return json.Unmarshal([]byte(s), v)
}

func JsonStringToMap[K comparable, V any, M ~map[K]V](s string) (M, error) {
	m := make(M)
	err := json.Unmarshal([]byte(s), &m)
	return m, err
}

func JsonStringToArray[E any, S ~[]E](s string) (S, error) {
	slice := make(S, 0)
	err := json.Unmarshal([]byte(s), &slice)
	return slice, err
}

func ToJsonString(a any) string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}
