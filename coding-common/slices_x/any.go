package slices_x

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
