package cmp_x

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func Test_com_00(t *testing.T) {
	testCases := []struct {
		old, new string
	}{
		// {"5.3.2", "5.3.2"},
		// {"v0.107.53", "v0.107.53"},
		// {"v0.107.53", "0.107.53"},
		{"12-alpine", "12.1-alpine"},
		// {"v27.0", "v27.1"},
		{"1.25.3-alpine3.19-perl", "1.25.4"},
		{"1.31.2", "1.31.1-alpine3.23-perl"},
		{"1.25.3-alpine3.22-perl", "1.25.3-alpine3.23-perl"},
		{"1.25.3-alpine22-perl", "1.25.3-alpine23-perl"},
		// {"9.0.0", "10.0.0"},
		{"1.0.0-alpha", "1.0.0"}, // alpha < beta < rc
		{"1.0.0-alpha", "1.0.0-beta"},
		{"1.0.0-rc.1", "1.0.0-rc.2"},
	}

	// a := strings.Compare
	a := func(a, b string) int { return CompareVersion(a, b, nil) }

	for _, tc := range testCases {
		result := a(tc.old, tc.new)
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

func Test_com_01(t *testing.T) {
	ss := []string{"", "a", "A", "9", " ", "	", "10"}
	slices.Sort(ss)

	fmt.Println(strings.Compare("9.0.0", "10.0.0"), ss)
}
