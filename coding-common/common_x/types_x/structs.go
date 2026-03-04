package types_x

type KV[T1 comparable, T2 any] = Pair[T1, T2]

// ------------------------------------------------------------------------------------------------------------------------

type Pair[T1, T2 any] struct {
	K T1 `json:"key"`
	V T2 `json:"value"`
}

func NewPair[T1, T2 any](k T1, v T2) Pair[T1, T2] { return Pair[T1, T2]{k, v} }

func (p *Pair[T1, T2]) Key() T1          { return p.K }
func (p *Pair[T1, T2]) Value() T2        { return p.V }
func (p *Pair[T1, T2]) Unpack() (T1, T2) { return p.K, p.V }

// ------------------------------------------------------------------------------------------------------------------------

type Tuple[T1, T2, T3 any] struct {
	A T1 `json:"a"`
	B T2 `json:"b"`
	C T3 `json:"c"`
}

func NewTuple[T1, T2, T3 any](a T1, b T2, c T3) Tuple[T1, T2, T3] { return Tuple[T1, T2, T3]{a, b, c} }

func (t *Tuple[T1, T2, T3]) A_() T1               { return t.A }
func (t *Tuple[T1, T2, T3]) B_() T2               { return t.B }
func (t *Tuple[T1, T2, T3]) C_() T3               { return t.C }
func (t *Tuple[T1, T2, T3]) Unpack() (T1, T2, T3) { return t.A, t.B, t.C }

// ------------------------------------------------------------------------------------------------------------------------

type ListNode[T any] struct {
	Val       T
	Pre, Next *ListNode[T]
}

func NewListNode[T any](v T, pre, next *ListNode[T]) *ListNode[T] {
	return &ListNode[T]{Pre: pre, Val: v, Next: next}
}

// ------------------------------------------------------------------------------------------------------------------------
