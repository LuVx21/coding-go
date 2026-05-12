package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
)

func compareVersions(oldVer, newVer string) (int, error) {
	// 清理版本号：移除常见的 'v' 前缀
	// oldVer, newVer = strings.TrimPrefix(oldVer, "v"), strings.TrimPrefix(newVer, "v")
	v1, err := version.NewVersion(oldVer)
	if err != nil {
		return 0, fmt.Errorf("解析旧版本失败 %s: %w", oldVer, err)
	}

	v2, err := version.NewVersion(newVer)
	if err != nil {
		return 0, fmt.Errorf("解析新版本失败 %s: %w", newVer, err)
	}

	return v1.Compare(v2), nil
}
func Test_version_00(t *testing.T) {
	testCases := []struct {
		old, new string
	}{
		{"5.3.2", "5.3.2"},
		{"v0.107.53", "v0.107.53"},
		{"v0.107.53", "0.107.53"},
		{"12-alpine", "12.1-alpine"},
		{"v27.0", "v27.1"},
		{"1.25.3-alpine3.19-perl", "1.25.4"},
		{"1.31.0", "1.25.3-bullseye"},
		{"1.25.3-alpine3.19-otel", "1.25.3-alpine3.21-otel"},
		{"8-alpine", "32bit"},
		{"1", "1.1.0"},
	}

	for _, tc := range testCases {
		result, err := compareVersions(tc.old, tc.new)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			continue
		}

		switch {
		case result < 0:
			fmt.Printf("%s < %s\n", tc.old, tc.new)
		case result > 0:
			fmt.Printf("%s > %s\n", tc.old, tc.new)
		default:
			fmt.Printf("%s == %s\n", tc.old, tc.new)
		}
	}
}
