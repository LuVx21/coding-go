package a

import "github.com/luvx21/coding-go/coding-common/common_x/t"

type EmptyHolder = struct{}

type (
	M[K comparable, V any] = map[K]V
	CAM[K comparable]      = map[K]any
	AAM                    = map[any]any
	SAM                    = map[string]any
)

type (
	S[T any]         = []T
	NS[N t.Number]   = []N
	CS[C comparable] = []C
	AS               = []any
	SS               = []string
	ASS              = [][]any
	SAMS             = []SAM
)

type (
	JsonObject = SAM
	JsonArray  = AS
)

type (
	Row                     = SAM
	Rows                    = []Row
	Table[RowId comparable] = map[RowId]Row
)

type (
	Seq0              = func(yield func() bool)
	Comparator[T any] = func(x, y T) int
)

type (
	Function[T, R any] = func(T) R
	Consumer[T any]    = func(T)
	Supplier[T any]    = func() T
	Runnable           = func()
	Predicate[T any]   = Function[T, bool]

	BiFunction[I1, I2, R any] = func(I1, I2) R
	BiConsumer[I1, I2 any]    = func(I1, I2)
	BiPredicate[I1, I2 any]   = BiFunction[I1, I2, bool]
)
