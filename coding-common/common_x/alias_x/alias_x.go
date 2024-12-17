package alias_x

type EmptyHolder = struct{}

type MapAny2Any = map[any]any
type MapComparable2Any[K comparable] = map[K]any
type MapStr2Any = map[string]any

type Slice[T any] = []T
type SliceAny = []any
type SliceStr = []string
type SliceSlice = []SliceAny
type SliceMapStr2Any = []MapStr2Any
type SliceComparable[T comparable] = []T

type JsonObject = MapStr2Any
type JsonArray = SliceAny
type Row = SliceAny

type Seq0 = func(yield func() bool)
