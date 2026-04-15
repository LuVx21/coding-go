package cast_x

import "time"

func ToBool(i any) bool {
	v, _ := ToBoolE(i)
	return v
}

func ToTime(i any) time.Time {
	return ToTimeInDefaultLocation(i, time.UTC)
}

func ToTimeInDefaultLocation(i any, location *time.Location) time.Time {
	v, _ := ToTimeInDefaultLocationE(i, location)
	return v
}

func ToDuration(i any) time.Duration {
	v, _ := ToDurationE(i)
	return v
}

func ToFloat64(i any) float64 {
	v, _ := ToFloat64E(i)
	return v
}

func ToFloat32(i any) float32 {
	v, _ := ToFloat32E(i)
	return v
}

func ToInt64(i any) int64 {
	v, _ := ToInt64E(i)
	return v
}

func ToInt32(i any) int32 {
	v, _ := ToInt32E(i)
	return v
}

func ToInt16(i any) int16 {
	v, _ := ToInt16E(i)
	return v
}

func ToInt8(i any) int8 {
	v, _ := ToInt8E(i)
	return v
}

func ToInt(i any) int {
	v, _ := ToIntE(i)
	return v
}

func ToUint(i any) uint {
	v, _ := ToUintE(i)
	return v
}

func ToUint64(i any) uint64 {
	v, _ := ToUint64E(i)
	return v
}

func ToUint32(i any) uint32 {
	v, _ := ToUint32E(i)
	return v
}

func ToUint16(i any) uint16 {
	v, _ := ToUint16E(i)
	return v
}

func ToUint8(i any) uint8 {
	v, _ := ToUint8E(i)
	return v
}

func ToString(i any) string {
	v, _ := ToStringE(i)
	return v
}

func ToStringMapString(i any) map[string]string {
	v, _ := ToStringMapStringE(i)
	return v
}

func ToStringMapStringSlice(i any) map[string][]string {
	v, _ := ToStringMapStringSliceE(i)
	return v
}

func ToStringMapBool(i any) map[string]bool {
	v, _ := ToStringMapBoolE(i)
	return v
}

func ToStringMapInt(i any) map[string]int {
	v, _ := ToStringMapIntE(i)
	return v
}

func ToStringMapInt64(i any) map[string]int64 {
	v, _ := ToStringMapInt64E(i)
	return v
}

func ToStringMap(i any) map[string]any {
	v, _ := ToStringMapE(i)
	return v
}

func ToSlice(i any) []any {
	v, _ := ToSliceE(i)
	return v
}

func ToBoolSlice(i any) []bool {
	v, _ := ToBoolSliceE(i)
	return v
}

func ToStringSlice(i any) []string {
	v, _ := ToStringSliceE(i)
	return v
}

func ToIntSlice(i any) []int {
	v, _ := ToIntSliceE(i)
	return v
}

func ToDurationSlice(i any) []time.Duration {
	v, _ := ToDurationSliceE(i)
	return v
}
