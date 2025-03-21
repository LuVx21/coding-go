package test

import (
	"errors"
	"runtime"
	"strings"
	"testing"
)

type Step struct {
	Name     string // 测试用例名称
	Input    any    // 入参
	Expected any    // 预期输出
	IsErr    bool   // 异常
	err      error  // 异常类型
}

func OneOne(t *testing.T, tests []Step, fn func(v any) any) {
	OneOneError(t, tests, func(v any) (any, error) {
		return fn(v), nil
	})
}

func OneOneError(t *testing.T, tests []Step, fn func(v any) (any, error)) {
	for _, tt := range tests {
		name, input, expected, err := tt.Name, tt.Input, tt.Expected, tt.err
		if len(name) == 0 {
			name = printCallerName()
		}
		t.Run(name, func(t *testing.T) {
			result, resultErr := fn(input)
			if tt.IsErr {
				if !errors.Is(resultErr, err) {
					t.Errorf("测试用例不通过(异常类型不匹配)-> \n入参: %v \n实际: %v \n预期: %v", input, resultErr, err)
				}
				return
			}
			if result != expected {
				t.Errorf("测试用例不通过-> \n入参: %v \n实际: %v \n预期: %v", input, result, expected)
			}
		})
	}
}

func printCallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	name := runtime.FuncForPC(pc).Name()
	index := strings.LastIndex(name, "/") + 1
	return name[index:]
}
