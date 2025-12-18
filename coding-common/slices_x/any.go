package slices_x

// ToAnySliceE 入参类型一致
func ToAnySliceE[E any](s ...E) []any {
	f := func(a E) any { return a }
	return Transfer(f, s...)
}

// ToAnySlice 入参类型可随意
func ToAnySlice(s ...any) []any {
	return ToAnySliceE(s...)
}

func SliceAny[S ~[]E, E any](s S) []any {
	r := make([]any, len(s))
	for i := range s {
		r[i] = s[i]
	}
	return r
}

func TypedSlice[S ~[]any, E any](s S) []E {
	r := make([]E, 0, len(s))
	for i := range s {
		if v, ok := s[i].(E); ok {
			r = append(r, v)
		}
	}
	return r
}
