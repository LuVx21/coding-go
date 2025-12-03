package mongodb

import "go.mongodb.org/mongo-driver/bson"

// ArgsE Map切片
func ArgsM(keys []string, values ...any) (r []bson.M, ok bool) {
	if len(keys) != len(values) {
		return nil, false
	}
	for i := range keys {
		r = append(r, bson.M{keys[i]: values[i]})
	}
	return
}

// ArgsE 结构体切片
func ArgsE(keys []string, values ...any) (r bson.D, ok bool) {
	if len(keys) != len(values) {
		return nil, false
	}
	for i := range keys {
		r = append(r, bson.E{Key: keys[i], Value: values[i]})
	}
	return
}
