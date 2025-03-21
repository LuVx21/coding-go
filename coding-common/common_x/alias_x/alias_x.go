package alias_x

import "github.com/luvx21/coding-go/coding-common/common_x/types_x"

type EmptyHolder = struct{}

type Map[K comparable, V any] = map[K]V
type MapAny2Any = map[any]any
type MapComparable2Any[K comparable] = map[K]any
type MapStr2Any = map[string]any

type Slice[T any] = []T
type NumberSlice[T types_x.Number] = []T
type SliceAny = []any
type SliceStr = []string
type SliceSlice = []SliceAny
type SliceMapStr2Any = []MapStr2Any
type SliceComparable[T comparable] = []T

type JsonObject = MapStr2Any
type JsonArray = SliceAny
type Row = SliceAny

type Seq0 = func(yield func() bool)
type Comparator[T any] = func(x, y T) int
