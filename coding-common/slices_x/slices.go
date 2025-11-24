package slices_x

import (
	"math/rand"
	"reflect"
	"sort"

	"slices"

	"github.com/luvx21/coding-go/coding-common/common_x/funcs"
	"github.com/luvx21/coding-go/coding-common/common_x/types_x"
	"golang.org/x/exp/constraints"
)

// Contains 适合较多数据或重复检查存在的场景(空间换时间)
func Contains[S ~[]E, E comparable](s S) func(E) bool {
	m := make(map[E]struct{}, len(s))
	for i := range s {
		m[s[i]] = struct{}{}
	}
	return func(e E) bool {
		_, ok := m[e]
		return ok
	}
}

func Partition[S ~[]E, E any](s S, size int) []S {
	r := make([]S, 0)
	_len := len(s)
	if _len == 0 {
		return r
	}

	step, mod := _len/size, _len%size
	for i := range step {
		r = append(r, s[i*size:(i+1)*size])
	}
	if mod > 0 {
		r = append(r, s[_len-mod:_len])
	}
	return r
}

func First[S ~[]E, E any](s S) (E, bool) {
	if len(s) == 0 {
		var zero E
		return zero, false
	}
	return s[0], true
}

func Last[S ~[]E, E any](s S) (E, bool) {
	l := len(s)
	if l == 0 {
		var zero E
		return zero, false
	}
	return s[l-1], true
}

func FirstOr[S ~[]E, E any](s S, def E) E {
	if len(s) == 0 {
		return def
	}
	return s[0]
}

func LastOr[S ~[]E, E any](s S, def E) E {
	l := len(s)
	if l == 0 {
		return def
	}
	return s[l-1]
}

func Transfer[I, O any](f funcs.Function[I, O], s ...I) []O {
	r := make([]O, len(s))
	for i := range s {
		r[i] = f(s[i])
	}
	return r
}

// ToAnySliceE 入参类型一致
func ToAnySliceE[E any](s ...E) []any {
	f := func(a E) any { return a }
	return Transfer(f, s...)
}

func FilterTransfer[I, O any](filter funcs.Predicate[I], f funcs.Function[I, O], s ...I) []O {
	r := make([]O, 0)
	for i := range s {
		if filter(s[i]) {
			r = append(r, f(s[i]))
		}
	}
	return r
}

// ToAnySlice 入参类型可随意
func ToAnySlice(s ...any) []any {
	return ToAnySliceE(s...)
}

func IsEmpty[S ~[]E, E any](s S) (bool, S) {
	if len(s) == 0 {
		return true, s
	}
	return false, s
}

func ClearZeroRef[S ~[]E, E any](s S) S {
	r := make(S, 0, len(s))
	for i := range s {
		if !reflect.ValueOf(&s[i]).Elem().IsZero() {
			r = append(r, s[i])
		}
	}
	return r
}

// Intersect 交集
func Intersect[S ~[]E, E comparable](a, b S) S {
	var r S
	seen := make(map[E]struct{}, len(a))
	for i := range a {
		seen[a[i]] = struct{}{}
	}
	for i := range b {
		if _, ok := seen[b[i]]; ok {
			r = append(r, b[i])
			delete(seen, b[i])
		}
	}
	return r
}

// Diff 差集a-b
func Diff[S ~[]E, E comparable](a, b S) S {
	var r S
	mp := make(map[E]struct{}, len(b))
	for i := range b {
		mp[b[i]] = struct{}{}
	}
	for i := range a {
		if _, ok := mp[a[i]]; !ok {
			r = append(r, a[i])
		}
	}
	return r
}

// IsUnique 切片是否存在重复值
func IsUnique[S ~[]E, E comparable](s S) bool {
	m := make(map[E]struct{}, len(s))
	for i := range s {
		if _, ok := m[s[i]]; ok {
			return false
		}
		m[s[i]] = struct{}{}
	}
	return true
}

// Unique 切片去重实现
func Unique[S ~[]E, E comparable](arr S) S {
	l := len(arr)
	if l <= 1 {
		return arr
	}

	r := make(S, 0, l)
	mp := map[E]struct{}{}
	for i := range arr {
		if _, ok := mp[arr[i]]; !ok {
			mp[arr[i]] = struct{}{}
			r = append(r, arr[i])
		}
	}
	return r
}

func Delete[S ~[]E, E comparable](arr S, t E) S {
	r := make(S, 0)
	for i := range arr {
		if arr[i] != t {
			r = append(r, arr[i])
		}
	}
	return r
}

func AllTrue[S ~[]E, E any](s S, fn funcs.Predicate[E]) bool {
	for i := range s {
		if !fn(s[i]) {
			return false
		}
	}
	return true
}
func AnyTrue[S ~[]E, E any](s S, fn funcs.Predicate[E]) bool {
	return slices.ContainsFunc(s, fn)
}

func AllIn[S ~[]E, E comparable](values S, safelist ...E) bool {
	for i := range values {
		if !slices.Contains(safelist, values[i]) {
			return false
		}
	}
	return true
}

func IsSorted[S ~[]E, E constraints.Ordered](s S) bool {
	return sort.SliceIsSorted(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}
func Sort[S ~[]E, E constraints.Ordered](s S) S {
	return SortBy(s, func(i, j E) bool { return i < j })
}

// SortBy 新slice
func SortBy[S ~[]E, E any](s S, less funcs.BiPredicate[E, E]) S {
	if len(s) <= 1 {
		return s
	}
	r := make(S, len(s))
	copy(r, s)
	sort.Slice(r, func(i, j int) bool { return less(r[i], r[j]) })
	return r
}

func FirstN[S ~[]E, E any](ss S, n int) (r S) {
	for i := 0; i < len(ss) && n > 0; i++ {
		r = append(r, ss[i])
		n--
	}
	return
}
func LastN[S ~[]E, E any](s S, n int) (r S) {
	if len(s) < n {
		n = len(s)
	}
	return s[len(s)-n:]
}

func Filter[S ~[]E, E any](s S, predicate funcs.Predicate[E]) (r S) {
	for i := range s {
		if predicate(s[i]) {
			r = append(r, s[i])
		}
	}
	return
}

// 双层切片
func Flat[S ~[]E, E any](s []S) (r S) {
	return FlatMap(s, func(s S) []E { return s })
}
func FlatMap[S ~[]E, E, O any](s S, transfer func(E) []O) (r []O) {
	for i := range s {
		r = append(r, transfer(s[i])...)
	}
	return r
}

func GroupBy[S ~[]E, E any, K comparable, V any](s S, keyMapper funcs.Function[E, K], valueMapper funcs.Function[E, V]) map[K][]V {
	groups := make(map[K][]V)
	for i := range s {
		key, value := keyMapper(s[i]), valueMapper(s[i])
		groups[key] = append(groups[key], value)
	}
	return groups
}

func Reduce[S ~[]E, E any, O any](ss S, reducer funcs.BiFunction[E, O, O]) (el O) {
	if len(ss) == 0 {
		return
	}
	for i := range ss {
		el = reducer(ss[i], el)
	}
	return
}

func Zip[E1, E2 any](ss1 []E1, ss2 []E2) []types_x.Pair[E1, E2] {
	minLen := min(len(ss1), len(ss2))
	var r []types_x.Pair[E1, E2]
	for i := range minLen {
		r = append(r, types_x.NewPair(ss1[i], ss2[i]))
	}
	return r
}

func Sum[S ~[]E, E types_x.Number](s S) (sum E) {
	for i := range s {
		sum += s[i]
	}
	return
}

func Shuffle[S ~[]E, E any](ss S, source rand.Source) S {
	n := len(ss)
	if n <= 1 {
		return ss
	}

	shuffled := make(S, n)
	copy(shuffled, ss)

	rnd := rand.New(source)
	rnd.Shuffle(n, func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}

func ForEach[S ~[]E, E any](s S, f func(index int, item E)) {
	for i := range s {
		f(i, s[i])
	}
}
func ForEachUntil[S ~[]E, E any](s S, until func(index int, item E) bool) {
	for i := range s {
		if until(i, s[i]) {
			break
		}
	}
}
