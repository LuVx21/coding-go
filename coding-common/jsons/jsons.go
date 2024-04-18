package jsons

import "encoding/json"

func JsonStringToObject(s string, v interface{}) error {
    return json.Unmarshal([]byte(s), v)
}

func JsonStringToMap[K comparable, V any, M ~map[K]V](s string) (*M, error) {
    m := make(M)
    err := json.Unmarshal([]byte(s), &m)
    return &m, err
}

func JsonStringToArray[E any, S ~[]E](s string) (*S, error) {
    slice := make(S, 0)
    err := json.Unmarshal([]byte(s), &slice)
    return &slice, err
}
