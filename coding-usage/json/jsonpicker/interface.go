package jsonpicker

type IJsonValue interface {
	Get(key string) IJsonValue
	DeepGet(key string) IJsonValue
	ContainsKey(key string) bool

	GetAt(index int) IJsonValue

	Contains(v any) bool

	Value() any

	IsMap() bool
	IsSliceOrArray() bool
}
