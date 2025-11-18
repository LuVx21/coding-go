package syncx

import "testing"

func Test_sync_map_00(t *testing.T) {
	m := NewSyncMap[string, int]()
	m.Store("hello", 123)
	println(m.Load("hello"))
}
