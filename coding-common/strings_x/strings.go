package strings_x

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/luvx21/coding-go/coding-common/cmp_x"
)

func IsDigit(r rune) bool         { return cmp_x.BetweenAnd(r, '0', '9') }
func IsLowerLetter(r rune) bool   { return cmp_x.BetweenAnd(r, 'a', 'z') }
func IsUpperLetter(r rune) bool   { return cmp_x.BetweenAnd(r, 'A', 'Z') }
func IsLetter(r rune) bool        { return IsLowerLetter(r) || IsUpperLetter(r) }
func IsLetterOrDigit(r rune) bool { return IsDigit(r) || IsLetter(r) }

func IsBlank[T ~string](cs T) bool {
	if len(cs) == 0 {
		return true
	}
	for _, c := range cs {
		if !unicode.IsSpace(c) {
			return false
		}
	}
	return true
}

func IsEmpty[T ~string](cs T) bool {
	return len(cs) == 0
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MinRunes(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func MaxRunes(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Map[T any](s string, mapping func(string) T) T {
	return mapping(s)
}

func FirstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}
