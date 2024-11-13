package reflects

import (
    "fmt"
    "reflect"
)

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

func RealType(a interface{}) reflect.Kind {
    if t := reflect.TypeOf(a); t.Kind() != reflect.Pointer {
        return t.Kind()
    } else {
        return t.Elem().Kind()
    }
}

func Indirect(a interface{}) interface{} {
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

func IndirectToStringerOrError(a interface{}) interface{} {
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
