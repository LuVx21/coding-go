package strings_x

import (
    "regexp"
    "strings"
    "unicode"
    "unicode/utf8"
)

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
