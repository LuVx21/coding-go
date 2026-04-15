package t

import "golang.org/x/exp/constraints"

type (
	// Number 数字
	Number interface {
		constraints.Integer | constraints.Float
	}
)

type (
	// Pair 键值对
	Pair[T1, T2 any] struct {
		K T1 `json:"key"`
		V T2 `json:"value"`
	}
	// Tuple 三元组
	Tuple[T1, T2, T3 any] struct {
		A T1 `json:"a"`
		B T2 `json:"b"`
		C T3 `json:"c"`
	}
	// ListNode 链表节点
	ListNode[T any] struct {
		Val       T
		Pre, Next *ListNode[T]
	}
)
type (
	// KV 键值对
	KV[T1 comparable, T2 any]  = Pair[T1, T2]
	KVS[T1 comparable, T2 any] = []Pair[T1, T2]
)

// ------------------------------------------------------------------------------------------------------------------------

func KvsToMap[M ~map[K]V, K comparable, V any](kvs KVS[K, V]) M {
	r := make(map[K]V, len(kvs))
	for _, kv := range kvs {
		r[kv.K] = kv.V
	}
	return r
}

func MapToKvs[M ~map[K]V, K comparable, V any](m M) KVS[K, V] {
	r := make(KVS[K, V], 0)
	for k, v := range m {
		r = append(r, KV[K, V]{K: k, V: v})
	}
	return r
}

// ------------------------------------------------------------------------------------------------------------------------

func NewPair[T1, T2 any](k T1, v T2) Pair[T1, T2] { return Pair[T1, T2]{k, v} }

func (p *Pair[T1, T2]) Key() T1          { return p.K }
func (p *Pair[T1, T2]) Value() T2        { return p.V }
func (p *Pair[T1, T2]) Unpack() (T1, T2) { return p.K, p.V }

// ------------------------------------------------------------------------------------------------------------------------

func NewTuple[T1, T2, T3 any](a T1, b T2, c T3) Tuple[T1, T2, T3] { return Tuple[T1, T2, T3]{a, b, c} }

func (t *Tuple[T1, T2, T3]) A_() T1               { return t.A }
func (t *Tuple[T1, T2, T3]) B_() T2               { return t.B }
func (t *Tuple[T1, T2, T3]) C_() T3               { return t.C }
func (t *Tuple[T1, T2, T3]) Unpack() (T1, T2, T3) { return t.A, t.B, t.C }

// ------------------------------------------------------------------------------------------------------------------------

func NewListNode[T any](v T, pre, next *ListNode[T]) *ListNode[T] {
	return &ListNode[T]{Pre: pre, Val: v, Next: next}
}
func (n *ListNode[T]) Data() T             { return n.Val }
func (n *ListNode[T]) Prev() *ListNode[T]  { return n.Pre }
func (n *ListNode[T]) NextN() *ListNode[T] { return n.Next }

// ------------------------------------------------------------------------------------------------------------------------
