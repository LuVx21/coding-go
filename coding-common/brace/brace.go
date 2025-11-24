package brace

import (
	"strconv"
	"strings"

	"github.com/luvx21/coding-go/coding-common/strings_x"
)

// Expand 对输入字符串执行花括号展开，返回展开后的所有字符串（按出现顺序）。
// 支持嵌套、大括号内逗号分隔项、范围 a..z 或 1..10（可带步长）以及反斜杠转义。
func Expand(s string) []string {
	// 查找第一个未被转义的 '{'
	i := 0
	for i < len(s) {
		if s[i] == '\\' {
			i += 2
			continue
		}
		if s[i] == '{' {
			break
		}
		i++
	}
	if i >= len(s) {
		// 无大括号，返回去转义后的原始字符串
		return []string{unescape(s)}
	}

	// 找到匹配的 '}'（考虑嵌套和转义）
	j := findMatchingBrace(s, i)
	if j == -1 {
		// 没有匹配的右花括号，视为字面字符串
		return []string{unescape(s)}
	}

	prefix := unescape(s[:i])
	content := s[i+1 : j]
	suffixPart := s[j+1:]

	// 解析大括号内部为选项（如果是序列则返回序列元素）
	options := parseBraceContent(content)

	// 对每个选项递归展开（因为选项内部可能还有大括号）
	var results []string
	for _, opt := range options {
		optExpanded, suffixExpanded := Expand(opt), Expand(suffixPart)
		for _, oe := range optExpanded {
			for _, se := range suffixExpanded {
				results = append(results, prefix+oe+se)
			}
		}
	}
	return results
}

// unescape 去掉反斜杠转义（仅移除反斜杠的含义，保留后续字符）
func unescape(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			b.WriteByte(s[i+1])
			i++
		} else {
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

// findMatchingBrace 从 idx (指向 '{') 开始寻找对应的 '}'，返回索引或 -1。
func findMatchingBrace(s string, idx int) int {
	if idx < 0 || idx >= len(s) || s[idx] != '{' {
		return -1
	}
	depth := 0
	for i := idx; i < len(s); i++ {
		if s[i] == '\\' {
			i++ // 跳过转义字符后面的字符
			continue
		}
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// parseBraceContent 将大括号内部解析为选项列表（逗号分隔）或序列元素。
func parseBraceContent(content string) []string {
	// 先判断是否为序列：形如 N..M 或 a..z，可选第三个参数步长
	if seq, ok := tryParseSequence(content); ok {
		return seq
	}
	// 否则按顶层逗号分割（考虑嵌套和转义）
	return splitTopLevel(content, ',')
}

// splitTopLevel 在顶层（不进入嵌套大括号）按 sep 分割字符串，返回包含空字符串的切片（与 shell 保持一致）。
func splitTopLevel(s string, sep rune) []string {
	var parts []string
	var cur strings.Builder
	depth, escaped := 0, false
	for _, r := range s {
		if escaped {
			cur.WriteRune(r)
			escaped = false
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		switch r {
		case '{':
			depth++
		case '}':
			if depth > 0 {
				depth--
			}
		}
		if r == sep && depth == 0 {
			parts = append(parts, cur.String())
			cur.Reset()
			continue
		}
		cur.WriteRune(r)
	}
	parts = append(parts, cur.String())
	return parts
}

// tryParseSequence 识别是否为序列形式并返回展开的元素
// 支持：整数序列（含负数）、字母序列（a..z）、可选步长（..step）
// 返回 (elements, true) 或 (nil, false)
func tryParseSequence(s string) ([]string, bool) {
	// 首先将未转义的点连续看成 .. delimiters，注意转义
	// 这里简单实现：在顶层无其他嵌套的情况下，用 splitTopLevel 按 '.' 分割可能比较麻烦
	// 使用手动解析识别形如 X..Y 或 X..Y..STEP（且没有未转义的逗号或花括号）
	// 如果内部包含未转义的逗号或花括号，则认为不是序列
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' {
			i++
			continue
		}
		if s[i] == '{' || s[i] == '}' || s[i] == ',' {
			return nil, false
		}
	}

	// 找到未转义的 ".." 的位置（第一个）
	pos := -1
	for i := 0; i+1 < len(s); i++ {
		if s[i] == '\\' {
			i++
			continue
		}
		if s[i] == '.' && s[i+1] == '.' {
			pos = i
			break
		}
	}
	if pos == -1 {
		return nil, false
	}
	// 再查找第二个 ".."（用于步长）
	pos2 := -1
	for i := pos + 2; i+1 < len(s); i++ {
		if s[i] == '\\' {
			i++
			continue
		}
		if s[i] == '.' && s[i+1] == '.' {
			pos2 = i
			break
		}
	}

	var aStr, bStr, stepStr string
	if pos2 == -1 {
		aStr, bStr, stepStr = s[:pos], s[pos+2:], ""
	} else {
		aStr, bStr, stepStr = s[:pos], s[pos+2:pos2], s[pos2+2:]
	}

	aStr, bStr, stepStr = strings.TrimSpace(aStr), strings.TrimSpace(bStr), strings.TrimSpace(stepStr)

	// 尝试作为整数序列
	if aInt, err1 := strconv.Atoi(aStr); err1 == nil {
		if bInt, err2 := strconv.Atoi(bStr); err2 == nil {
			step := 1
			if stepStr != "" {
				if stepVal, err3 := strconv.Atoi(stepStr); err3 == nil && stepVal != 0 {
					step = stepVal
				} else {
					return nil, false
				}
			}
			return expandIntSequence(aInt, bInt, step), true
		}
		// not both ints
		return nil, false
	}

	// 尝试作为字母序列（单字符）
	if len([]rune(aStr)) == 1 && len([]rune(bStr)) == 1 {
		ra, rb := []rune(aStr)[0], []rune(bStr)[0]
		// 仅支持 ASCII 字母/数字范围中的单字符序列
		if strings_x.IsLetterOrDigit(ra) && strings_x.IsLetterOrDigit(rb) {
			step := 1
			if stepStr != "" {
				if stepVal, err := strconv.Atoi(stepStr); err == nil && stepVal != 0 {
					step = stepVal
				} else {
					return nil, false
				}
			}
			return expandRuneSequence(ra, rb, step), true
		}
	}

	return nil, false
}

func expandIntSequence(a, b, step int) []string {
	var res []string
	if step == 0 {
		return res
	}
	up := a <= b
	if up {
		if step < 0 {
			step = -step
		}
		for i := a; i <= b; i += step {
			res = append(res, strconv.Itoa(i))
		}
	} else {
		if step > 0 {
			step = -step
		}
		for i := a; i >= b; i += step {
			res = append(res, strconv.Itoa(i))
		}
	}
	return res
}

func expandRuneSequence(a, b rune, step int) []string {
	var res []string
	if step == 0 {
		return res
	}
	up := a <= b
	if up {
		if step < 0 {
			step = -step
		}
		for i := a; i <= b; i += rune(step) {
			res = append(res, string(i))
		}
	} else {
		if step > 0 {
			step = -step
		}
		for i := a; i >= b; i += rune(step) {
			res = append(res, string(i))
		}
	}
	return res
}
