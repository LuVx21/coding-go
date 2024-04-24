package strings_x

import "unicode"

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
