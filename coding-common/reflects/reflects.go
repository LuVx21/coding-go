package reflects

import (
	"fmt"
	"reflect"

	"github.com/luvx21/coding-go/coding-common/slices_x"
)

func IsZeroRef[T any](v T) bool {
	return reflect.ValueOf(&v).Elem().IsZero()
}

func IsNil(v any) bool {
	valueOf := reflect.ValueOf(&v).Elem()
	k := valueOf.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return valueOf.IsNil()
	default:
		return v == nil
	}
}

func RealType(a any) reflect.Kind {
	if t := reflect.TypeOf(a); t.Kind() != reflect.Pointer {
		return t.Kind()
	} else {
		return t.Elem().Kind()
	}
}

func Indirect(a any) any {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Pointer {
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

func IndirectToStringerOrError(a any) any {
	if a == nil {
		return nil
	}

	errorType := reflect.TypeOf((*error)(nil)).Elem()
	fmtStringerType := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

	v := reflect.ValueOf(a)
	for !v.Type().Implements(fmtStringerType) && !v.Type().Implements(errorType) && v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// CallFunc 反射调用函数
func CallFunc(runFn any, args ...any) {
	fn := reflect.ValueOf(runFn)

	fnType := fn.Type()
	inNum := fnType.NumIn()

	argArray := slices_x.Partition(args, inNum)
	for _, arg := range argArray {
		args := slices_x.Transfer(func(i any) reflect.Value {
			return reflect.ValueOf(i)
		}, arg...)
		result := fn.Call(args)
		fmt.Printf("入参: %v -> 结果:%v\n", arg, result[0])
	}
}
