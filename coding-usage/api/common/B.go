package common

type MyInt int

func (m MyInt) IsZero() bool {
    return m == 0
}

func (m MyInt) Add(other int) int {
    return other + int(m)
}
