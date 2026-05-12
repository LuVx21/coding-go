package main

import (
	"fmt"
	"testing"
)

func Test_version_00(t *testing.T) {
	testCases := []struct{ tag string }{
		{"latest"},
		{"v4.0.3-ls215"},
		{"vscode"},
	}
	for _, c := range testCases {
		v, s, b := FromTag(c.tag)
		fmt.Println(c.tag, "->", b, v, s)
	}
}

func Test_version_01(t *testing.T) {
	testCases := []struct {
		old, new string
	}{
		{"12-alpine", "12.1-alpine"},
		{"1.25.3-alpine3.19-otel", "1.25.3-alpine3.21-otel"},
		{"1", "1.1.0"},

		{"5.3.2", "5.3.2"},
		{"v0.107.53", "v0.107.53"},
		{"v0.107.53", "0.107.53"},
		{"v27.0", "v27.1"},
		{"1.25.3-alpine3.19-perl", "1.25.4"},
		{"1.31.0", "1.25.3-bullseye"},
		{"8-alpine", "32bit"},
	}

	for _, tc := range testCases {
		old, _, _ := FromTag(tc.old)
		new_, _, _ := FromTag(tc.new)

		result := old.Compare(new_)

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
