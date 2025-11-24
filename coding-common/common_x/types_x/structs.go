package types_x

type KV[T1 comparable, T2 any] = Pair[T1, T2]

// type KV[K comparable, V any] struct {
// 	Key   K `json:"key"`
// 	Value V `json:"value"`
// }

type Pair[T1, T2 any] struct {
	K T1 `json:"key"`
	V T2 `json:"value"`
}

func NewPair[T1, T2 any](k T1, v T2) Pair[T1, T2] { return Pair[T1, T2]{k, v} }

func (p *Pair[T1, T2]) Key() T1          { return p.K }
func (p *Pair[T1, T2]) Value() T2        { return p.V }
func (p *Pair[T1, T2]) Unpack() (T1, T2) { return p.K, p.V }
