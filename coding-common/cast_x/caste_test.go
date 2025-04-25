package cast_x

import (
	"encoding/json"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStep struct {
	input  any
	expect any
	iserr  bool
}

func commonData(zero, one, expect any) []testStep {
	return []testStep{
		{int(121), expect, false},
		{int8(121), expect, false},
		{int16(121), expect, false},
		{int32(121), expect, false},
		{int64(121), expect, false},
		{uint(121), expect, false},
		{uint8(121), expect, false},
		{uint16(121), expect, false},
		{uint32(121), expect, false},
		{uint64(121), expect, false},
		{"121", expect, false},
		{json.Number("121"), expect, false},
		{float64(121.1), expect, false},
		{float32(121.1), expect, false},
		{true, one, false},
		{false, zero, false},
	}
}

func test(t *testing.T,
	dataSlice []testStep,
	f func(v any) (any, error),
) {
	name := printCallerName()
	for _, data := range dataSlice {
		input, expected := data.input, data.expect
		actual, err := f(input)
		if data.iserr {
			assert.Error(t, err, "方法:%s, 测试用例:%v", name, data)
			continue
		}
		assert.Equal(t, expected, actual, "方法:%s, 测试用例:%v", name, data)
	}
}

func printCallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	name := runtime.FuncForPC(pc).Name()
	index := strings.LastIndex(name, "/") + 1
	return name[index:]
}

func TestToInt64E(t *testing.T) {
	datas := commonData(int64(0), int64(1), int64(121))
	test(t, datas, func(v any) (any, error) { return ToInt64E(v) })
}

func TestToUint64E(t *testing.T) {
	datas := commonData(uint64(0), uint64(1), uint64(121))
	test(t, datas, func(v any) (any, error) { return ToUint64E(v) })
}
