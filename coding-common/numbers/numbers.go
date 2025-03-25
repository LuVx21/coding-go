package numbers

import (
	"strings"
)

// TrimZeroDecimal 保留小数点之前
func TrimZeroDecimal(s string) string {
	for i := range len(s) {
		if s[i] < '0' || s[i] > '9' {
			if s[i] == '.' {
				return s[0:i]
			}
			if i != 0 || (s[i] != '+' && s[i] != '-') {
				return s
			}
		}
	}
	return s
}

// TruncateDecimals 保留小数点之前
func TruncateDecimals(s string) string {
	v := strings.LastIndex(s, ".")
	if v == -1 {
		return s
	}

	// 小数点后全是数字,执行完整个循环不退出
	for i := len(s) - 1; i > v; i-- {
		if s[i] < '0' || s[i] > '9' {
			return s
		}
	}
	// 小数点之前,存在非数字,不处理
	for i := v - 1; i >= 0; i-- {
		if s[i] < '0' || s[i] > '9' {
			if i != 0 || (s[i] != '+' && s[i] != '-') {
				return s
			}
		}
	}
	return s[:v]
}
