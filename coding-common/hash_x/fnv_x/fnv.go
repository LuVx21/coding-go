package fnv_x

import "fmt"

const (
    offset32 = 2166136261
    prime32  = uint32(16777619)
)

func StrFnv32[K fmt.Stringer](key K) uint32 {
    return Fnv32(key.String())
}

func Fnv32(key string) uint32 {
    hash := uint32(offset32)
    keyLength := len(key)
    for i := range keyLength {
        hash *= prime32
        hash ^= uint32(key[i])
    }
    return hash
}
