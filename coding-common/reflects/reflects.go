package reflects

import "reflect"

func IsZeroRef[T any](v T) bool {
    return reflect.ValueOf(&v).Elem().IsZero()
}

func IsNil(v interface{}) bool {
    valueOf := reflect.ValueOf(&v).Elem()
    k := valueOf.Kind()
    switch k {
    case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
        return valueOf.IsNil()
    default:
        return v == nil
    }
}
