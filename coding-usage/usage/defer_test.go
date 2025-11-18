package main

import "testing"

func around(t *testing.T, msg string) func(t *testing.T) {
	t.Log(msg + ": 测试前")
	return func(t *testing.T) {
		t.Log(msg + ": 测试后")
	}
}

func Test_00(t *testing.T) {
	defer around(t, "Test_00")(t)

	t.Log("执行...")
}
