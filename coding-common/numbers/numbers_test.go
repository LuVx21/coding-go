package numbers

import (
    "fmt"
    "testing"
)

func Test_00(t *testing.T) {}

func Test_01(t *testing.T) {
    fmt.Println(
        TrimZeroDecimal("1.010"),
        TrimZeroDecimal("1.0a10"),
        TrimZeroDecimal("+1.010"),
        TrimZeroDecimal("+1a.010"),
        TrimZeroDecimal("+a1.010"),
        TrimZeroDecimal(".11"),
        TrimZeroDecimal("0.11"),
        TrimZeroDecimal("a1.11"),
    )
    decimals := TruncateDecimals("1.010")
    fmt.Println(decimals)
}
