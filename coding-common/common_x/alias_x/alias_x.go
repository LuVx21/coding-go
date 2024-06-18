package alias_x

type EmptyHolder = struct{}

type MapAny2Any = map[any]any
type MapStr2Any = map[string]any

type SliceAny = []any
type SliceStr = []string
type SliceSlice = []SliceAny
type SliceMapStr2Any = []MapStr2Any

type JsonObject = MapStr2Any
type JsonArray = SliceAny

type Seq0 = func(yield func() bool)
