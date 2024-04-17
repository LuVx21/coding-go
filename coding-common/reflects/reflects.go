package reflects

import "reflect"

func IsZeroRef[T any](v T) bool {
    return reflect.ValueOf(&v).Elem().IsZero()
}

func IsNil(v interface{}) bool {
    valueOf := reflect.ValueOf(&v).Elem()
    k := valueOf.Kind()
    switch k {
    case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
        return valueOf.IsNil()
    default:
        return v == nil
    }
}

func Indirect(a interface{}) interface{} {
    if a == nil {
        return nil
    }
    if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
        return a
    }
    v := reflect.ValueOf(a)
    for v.Kind() == reflect.Ptr && !v.IsNil() {
        v = v.Elem()
    }
    return v.Interface()
}
