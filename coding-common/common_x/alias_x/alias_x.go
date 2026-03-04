package alias_x

import "github.com/luvx21/coding-go/coding-common/common_x/types_x"

type EmptyHolder = struct{}

type (
	Map[K comparable, V any]        = map[K]V
	MapAny2Any                      = map[any]any
	MapComparable2Any[K comparable] = map[K]any
	MapStr2Any                      = map[string]any
)

type (
	Slice[T any]                  = []T
	NumberSlice[T types_x.Number] = []T
	SliceAny                      = []any
	SliceStr                      = []string
	SliceSlice                    = []SliceAny
	SliceMapStr2Any               = []MapStr2Any
	SliceComparable[T comparable] = []T
)

type JsonObject = MapStr2Any
type JsonArray = SliceAny

type Row = MapStr2Any
type Rows = []Row

type Seq0 = func(yield func() bool)
type Comparator[T any] = func(x, y T) int
