package cast_x

import "time"

func ToBool(i interface{}) bool {
	v, _ := ToBoolE(i)
	return v
}

func ToTime(i interface{}) time.Time {
	return ToTimeInDefaultLocation(i, time.UTC)
}

func ToTimeInDefaultLocation(i interface{}, location *time.Location) time.Time {
	v, _ := ToTimeInDefaultLocationE(i, location)
	return v
}

func ToDuration(i interface{}) time.Duration {
	v, _ := ToDurationE(i)
	return v
}

func ToFloat64(i interface{}) float64 {
	v, _ := ToFloat64E(i)
	return v
}

func ToFloat32(i interface{}) float32 {
	v, _ := ToFloat32E(i)
	return v
}

func ToInt64(i interface{}) int64 {
	v, _ := ToInt64E(i)
	return v
}

func ToInt32(i interface{}) int32 {
	v, _ := ToInt32E(i)
	return v
}

func ToInt16(i interface{}) int16 {
	v, _ := ToInt16E(i)
	return v
}

func ToInt8(i interface{}) int8 {
	v, _ := ToInt8E(i)
	return v
}

func ToInt(i interface{}) int {
	v, _ := ToIntE(i)
	return v
}

func ToUint(i interface{}) uint {
	v, _ := ToUintE(i)
	return v
}

func ToUint64(i interface{}) uint64 {
	v, _ := ToUint64E(i)
	return v
}

func ToUint32(i interface{}) uint32 {
	v, _ := ToUint32E(i)
	return v
}

func ToUint16(i interface{}) uint16 {
	v, _ := ToUint16E(i)
	return v
}

func ToUint8(i interface{}) uint8 {
	v, _ := ToUint8E(i)
	return v
}

func ToString(i interface{}) string {
	v, _ := ToStringE(i)
	return v
}

func ToStringMapString(i interface{}) map[string]string {
	v, _ := ToStringMapStringE(i)
	return v
}

func ToStringMapStringSlice(i interface{}) map[string][]string {
	v, _ := ToStringMapStringSliceE(i)
	return v
}

func ToStringMapBool(i interface{}) map[string]bool {
	v, _ := ToStringMapBoolE(i)
	return v
}

func ToStringMapInt(i interface{}) map[string]int {
	v, _ := ToStringMapIntE(i)
	return v
}

func ToStringMapInt64(i interface{}) map[string]int64 {
	v, _ := ToStringMapInt64E(i)
	return v
}

func ToStringMap(i interface{}) map[string]interface{} {
	v, _ := ToStringMapE(i)
	return v
}

func ToSlice(i interface{}) []interface{} {
	v, _ := ToSliceE(i)
	return v
}

func ToBoolSlice(i interface{}) []bool {
	v, _ := ToBoolSliceE(i)
	return v
}

func ToStringSlice(i interface{}) []string {
	v, _ := ToStringSliceE(i)
	return v
}

func ToIntSlice(i interface{}) []int {
	v, _ := ToIntSliceE(i)
	return v
}

func ToDurationSlice(i interface{}) []time.Duration {
	v, _ := ToDurationSliceE(i)
	return v
}
