package types_x

type KV[K comparable, V any] struct {
	Key   K `json:"key"`
	Value V `json:"value"`
}
