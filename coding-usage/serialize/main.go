package serialize

type Serializer[T any] interface {
	Marshal(T) ([]byte, error)
	Unmarshal([]byte, T) error
}
