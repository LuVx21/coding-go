package maths_x

import "golang.org/x/exp/constraints"

// 实现细节：先把值转换为 int64 做运算，再转换回 T。
// 注意：如果目标类型比 int64 更宽（在目前 Go 常见类型中不会），
// 或者计算导致溢出，结果会被截断；如需严格溢出检测需额外判断。
func Add[A, B, C constraints.Signed](a A, b B) C {
	return C(int64(a) + int64(b))
}
