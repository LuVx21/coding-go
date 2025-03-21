package slices_x

import (
	"math/rand"
	"reflect"
	"sort"

	. "github.com/luvx21/coding-go/coding-common/common_x/funcs"
	. "github.com/luvx21/coding-go/coding-common/common_x/pairs"
	. "github.com/luvx21/coding-go/coding-common/common_x/types_x"
	"golang.org/x/exp/constraints"
)

// Contains 适合较多数据或重复检查存在的场景(空间换时间)
func Contains[S ~[]E, E comparable](s S) func(E) bool {
	m := make(map[E]struct{})
	for _, e := range s {
		m[e] = struct{}{}
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

func Transfer[I, O any](f Function[I, O], s ...I) []O {
	r := make([]O, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

// ToAnySliceE 入参类型一致
func ToAnySliceE[E any](s ...E) []any {
	f := func(a E) any { return a }
	return Transfer[E, any](f, s...)
}

func FilterTransfer[I, O any](filter Predicate[I], f Function[I, O], s ...I) []O {
	r := make([]O, 0)
	for _, e := range s {
		if filter(e) {
			r = append(r, f(e))
		}
	}
	return r
}

// ToAnySlice 入参类型可随意
func ToAnySlice(s ...any) []any {
	return ToAnySliceE[any](s...)
}

func IsEmpty[S ~[]E, E any](s S) (bool, S) {
	if s == nil || len(s) == 0 {
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
	mp := make(map[E]struct{}, len(a))
	for _, val := range a {
		mp[val] = struct{}{}
	}
	for _, val := range b {
		if _, ok := mp[val]; ok {
			r = append(r, val)
		}
	}
	return r
}

// Diff 差集a-b
func Diff[S ~[]E, E comparable](a, b S) S {
	var r S
	mp := make(map[E]struct{}, len(b))
	for _, val := range b {
		mp[val] = struct{}{}
	}
	for _, val := range a {
		if _, ok := mp[val]; !ok {
			r = append(r, val)
		}
	}
	return r
}

// IsUnique 切片是否存在重复值
func IsUnique[S ~[]E, E comparable](s S) bool {
	m := make(map[E]struct{}, len(s))
	for _, e := range s {
		if _, ok := m[e]; ok {
			return false
		}
		m[e] = struct{}{}
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
	for _, e := range arr {
		if _, ok := mp[e]; !ok {
			mp[e] = struct{}{}
			r = append(r, e)
		}
	}
	return r
}

func Delete[S ~[]E, E comparable](arr S, t E) S {
	r := make(S, 0)
	for _, e := range arr {
		if e != t {
			r = append(r, e)
		}
	}
	return r
}

func AllTrue[S ~[]E, E any](s S, fn Predicate[E]) bool {
	for _, value := range s {
		if !fn(value) {
			return false
		}
	}
	return true
}
func AnyTrue[S ~[]E, E any](s S, fn Predicate[E]) bool {
	for _, value := range s {
		if fn(value) {
			return true
		}
	}
	return false
}

func In[S ~[]E, E comparable](value E, safelist ...E) bool {
	for i := range safelist {
		if value == safelist[i] {
			return true
		}
	}
	return false
}
func AllIn[S ~[]E, E comparable](values S, safelist ...E) bool {
	for i := range values {
		if !In[S, E](values[i], safelist...) {
			return false
		}
	}
	return true
}
func NotIn[S ~[]E, E comparable](value E, blocklist ...E) bool {
	for i := range blocklist {
		if value == blocklist[i] {
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
	return SortBy[S, E](s, func(i, j E) bool { return i < j })
}

// SortBy 新slice
func SortBy[S ~[]E, E any](s S, less BiPredicate[E, E]) S {
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
	var lastIndex = len(s) - 1
	for i := lastIndex; i >= 0 && n > 0; i-- {
		r = append(r, s[i])
		n--
	}
	return
}

func Filter[S ~[]E, E any](s S, predicate Predicate[E]) (r S) {
	for _, e := range s {
		if predicate(e) {
			r = append(r, e)
		}
	}
	return
}

func Flat[S ~[]E, E any](s []S) (r S) {
	for _, e := range s {
		r = append(r, e...)
	}
	return r
}

func GroupBy[S ~[]E, E any, K comparable](s S, by Function[E, K]) map[K]S {
	groups := make(map[K]S)
	for _, e := range s {
		key := by(e)
		groups[key] = append(groups[key], e)
	}
	return groups
}

func Reduce[S ~[]E, E any, O any](ss S, reducer BiFunction[E, O, O]) (el O) {
	if len(ss) == 0 {
		return
	}
	for _, s := range ss {
		el = reducer(s, el)
	}
	return
}

func Zip[E1, E2 any](ss1 []E1, ss2 []E2) []Pair[E1, E2] {
	minLen := min(len(ss1), len(ss2))
	var r []Pair[E1, E2]
	for i := 0; i < minLen; i++ {
		r = append(r, NewPair(ss1[i], ss2[i]))
	}
	return r
}

func Sum[S ~[]E, E Number](s S) (sum E) {
	for _, s := range s {
		sum += s
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
