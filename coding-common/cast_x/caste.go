package cast_x

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/luvx21/coding-go/coding-common/jsons"
	"github.com/luvx21/coding-go/coding-common/numbers"
	"github.com/luvx21/coding-go/coding-common/reflects"
	"github.com/luvx21/coding-go/coding-common/times_x"
)

var (
	errNegativeNotAllowed = errors.New("unable to cast negative value")
)

type float64EProvider interface {
	Float64() (float64, error)
}
type float64Provider interface {
	Float64() float64
}

func ToTimeInDefaultLocationE(i interface{}, location *time.Location) (tim time.Time, err error) {
	i = reflects.Indirect(i)

	switch rv := i.(type) {
	case time.Time:
		return rv, nil
	case string:
		return times_x.StringToDateInDefaultLocation(rv, location)
	case json.Number:
		s, err1 := ToInt64E(rv)
		if err1 != nil {
			return time.Time{}, fmt.Errorf("unable to cast %#v of type %T to Time", i, i)
		}
		return time.Unix(s, 0), nil
	case int:
		return time.Unix(int64(rv), 0), nil
	case int64:
		return time.Unix(rv, 0), nil
	case int32:
		return time.Unix(int64(rv), 0), nil
	case uint:
		return time.Unix(int64(rv), 0), nil
	case uint64:
		return time.Unix(int64(rv), 0), nil
	case uint32:
		return time.Unix(int64(rv), 0), nil
	default:
		return time.Time{}, fmt.Errorf("unable to cast %#v of type %T to Time", i, i)
	}
}

func ToDurationE(i interface{}) (d time.Duration, err error) {
	i = reflects.Indirect(i)

	switch rv := i.(type) {
	case time.Duration:
		return rv, nil
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		d = time.Duration(ToInt64(rv))
		return
	case float32, float64:
		d = time.Duration(ToFloat64(rv))
		return
	case string:
		if strings.ContainsAny(rv, "nsuµmh") {
			d, err = time.ParseDuration(rv)
		} else {
			d, err = time.ParseDuration(rv + "ns")
		}
		return
	case float64EProvider:
		var v float64
		v, err = rv.Float64()
		d = time.Duration(v)
		return
	case float64Provider:
		d = time.Duration(rv.Float64())
		return
	default:
		err = fmt.Errorf("unable to cast %#v of type %T to Duration", i, i)
		return
	}
}

func ToBoolE(i interface{}) (bool, error) {
	i = reflects.Indirect(i)

	switch rv := i.(type) {
	case bool:
		return rv, nil
	case nil:
		return false, nil
	case int:
		return rv != 0, nil
	case int64:
		return rv != 0, nil
	case int32:
		return rv != 0, nil
	case int16:
		return rv != 0, nil
	case int8:
		return rv != 0, nil
	case uint:
		return rv != 0, nil
	case uint64:
		return rv != 0, nil
	case uint32:
		return rv != 0, nil
	case uint16:
		return rv != 0, nil
	case uint8:
		return rv != 0, nil
	case float64:
		return rv != 0, nil
	case float32:
		return rv != 0, nil
	case time.Duration:
		return rv != 0, nil
	case string:
		return strconv.ParseBool(numbers.TrimZeroDecimal(rv))
	case json.Number:
		return strconv.ParseBool(numbers.TrimZeroDecimal(string(rv)))
	default:
		return false, fmt.Errorf("unable to cast %#v of type %T to bool", i, i)
	}
}

func ToFloat64E(i interface{}) (float64, error) {
	i = reflects.Indirect(i)

	intv, ok := toInt(i)
	if ok {
		return float64(intv), nil
	}

	switch rv := i.(type) {
	case float64:
		return rv, nil
	case float32:
		return float64(rv), nil
	case int64:
		return float64(rv), nil
	case int32:
		return float64(rv), nil
	case int16:
		return float64(rv), nil
	case int8:
		return float64(rv), nil
	case uint:
		return float64(rv), nil
	case uint64:
		return float64(rv), nil
	case uint32:
		return float64(rv), nil
	case uint16:
		return float64(rv), nil
	case uint8:
		return float64(rv), nil
	case string:
		v, err := strconv.ParseFloat(rv, 64)
		if err == nil {
			return v, nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to float64", i, i)
	case float64EProvider:
		v, err := rv.Float64()
		if err == nil {
			return v, nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to float64", i, i)
	case float64Provider:
		return rv.Float64(), nil
	case bool:
		if rv {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to float64", i, i)
	}
}

func ToFloat32E(i interface{}) (float32, error) {
	i = reflects.Indirect(i)

	intv, ok := toInt(i)
	if ok {
		return float32(intv), nil
	}

	switch rv := i.(type) {
	case float64:
		return float32(rv), nil
	case float32:
		return rv, nil
	case int64:
		return float32(rv), nil
	case int32:
		return float32(rv), nil
	case int16:
		return float32(rv), nil
	case int8:
		return float32(rv), nil
	case uint:
		return float32(rv), nil
	case uint64:
		return float32(rv), nil
	case uint32:
		return float32(rv), nil
	case uint16:
		return float32(rv), nil
	case uint8:
		return float32(rv), nil
	case string:
		v, err := strconv.ParseFloat(rv, 32)
		if err == nil {
			return float32(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to float32", i, i)
	case float64EProvider:
		v, err := rv.Float64()
		if err == nil {
			return float32(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to float32", i, i)
	case float64Provider:
		return float32(rv.Float64()), nil
	case bool:
		if rv {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to float32", i, i)
	}
}

func ToInt64E(i interface{}) (int64, error) {
	i = reflects.Indirect(i)

	if intv, ok := toInt(i); ok {
		return intv, nil
	}

	switch rv := i.(type) {
	case int64:
		return rv, nil
	case int32:
		return int64(rv), nil
	case int16:
		return int64(rv), nil
	case int8:
		return int64(rv), nil
	case uint:
		return int64(rv), nil
	case uint64:
		return int64(rv), nil
	case uint32:
		return int64(rv), nil
	case uint16:
		return int64(rv), nil
	case uint8:
		return int64(rv), nil
	case float64:
		return int64(rv), nil
	case float32:
		return int64(rv), nil
	case string:
		v, err := strconv.ParseInt(numbers.TrimZeroDecimal(rv), 0, 0)
		if err == nil {
			return v, nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
	case json.Number:
		return ToInt64E(string(rv))
	case bool:
		if rv {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
	}
}

func ToInt32E(i interface{}) (int32, error) {
	i = reflects.Indirect(i)

	intv, ok := toInt(i)
	if ok {
		return int32(intv), nil
	}

	switch rv := i.(type) {
	case int64:
		return int32(rv), nil
	case int32:
		return rv, nil
	case int16:
		return int32(rv), nil
	case int8:
		return int32(rv), nil
	case uint:
		return int32(rv), nil
	case uint64:
		return int32(rv), nil
	case uint32:
		return int32(rv), nil
	case uint16:
		return int32(rv), nil
	case uint8:
		return int32(rv), nil
	case float64:
		return int32(rv), nil
	case float32:
		return int32(rv), nil
	case string:
		v, err := strconv.ParseInt(numbers.TrimZeroDecimal(rv), 0, 0)
		if err == nil {
			return int32(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to int32", i, i)
	case json.Number:
		return ToInt32E(string(rv))
	case bool:
		if rv {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int32", i, i)
	}
}

func ToInt16E(i interface{}) (int16, error) {
	i = reflects.Indirect(i)

	intv, ok := toInt(i)
	if ok {
		return int16(intv), nil
	}

	switch rv := i.(type) {
	case int64:
		return int16(rv), nil
	case int32:
		return int16(rv), nil
	case int16:
		return rv, nil
	case int8:
		return int16(rv), nil
	case uint:
		return int16(rv), nil
	case uint64:
		return int16(rv), nil
	case uint32:
		return int16(rv), nil
	case uint16:
		return int16(rv), nil
	case uint8:
		return int16(rv), nil
	case float64:
		return int16(rv), nil
	case float32:
		return int16(rv), nil
	case string:
		v, err := strconv.ParseInt(numbers.TrimZeroDecimal(rv), 0, 0)
		if err == nil {
			return int16(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to int16", i, i)
	case json.Number:
		return ToInt16E(string(rv))
	case bool:
		if rv {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int16", i, i)
	}
}

func ToInt8E(i interface{}) (int8, error) {
	i = reflects.Indirect(i)

	intv, ok := toInt(i)
	if ok {
		return int8(intv), nil
	}

	switch rv := i.(type) {
	case int64:
		return int8(rv), nil
	case int32:
		return int8(rv), nil
	case int16:
		return int8(rv), nil
	case int8:
		return rv, nil
	case uint:
		return int8(rv), nil
	case uint64:
		return int8(rv), nil
	case uint32:
		return int8(rv), nil
	case uint16:
		return int8(rv), nil
	case uint8:
		return int8(rv), nil
	case float64:
		return int8(rv), nil
	case float32:
		return int8(rv), nil
	case string:
		v, err := strconv.ParseInt(numbers.TrimZeroDecimal(rv), 0, 0)
		if err == nil {
			return int8(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to int8", i, i)
	case json.Number:
		return ToInt8E(string(rv))
	case bool:
		if rv {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int8", i, i)
	}
}

func ToIntE(i interface{}) (int, error) {
	i = reflects.Indirect(i)

	intv, ok := toInt(i)
	if ok {
		return int(intv), nil
	}

	switch rv := i.(type) {
	case int64:
		return int(rv), nil
	case int32:
		return int(rv), nil
	case int16:
		return int(rv), nil
	case int8:
		return int(rv), nil
	case uint:
		return int(rv), nil
	case uint64:
		return int(rv), nil
	case uint32:
		return int(rv), nil
	case uint16:
		return int(rv), nil
	case uint8:
		return int(rv), nil
	case float64:
		return int(rv), nil
	case float32:
		return int(rv), nil
	case string:
		v, err := strconv.ParseInt(numbers.TrimZeroDecimal(rv), 0, 0)
		if err == nil {
			return int(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
	case json.Number:
		return ToIntE(string(rv))
	case bool:
		if rv {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int", i, i)
	}
}

func ToUintE(i interface{}) (uint, error) {
	i = reflects.Indirect(i)

	intv, ok := toInt(i)
	if ok {
		if intv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(intv), nil
	}

	switch rv := i.(type) {
	case string:
		v, err := strconv.ParseInt(numbers.TrimZeroDecimal(rv), 0, 0)
		if err == nil {
			if v < 0 {
				return 0, errNegativeNotAllowed
			}
			return uint(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint", i, i)
	case json.Number:
		return ToUintE(string(rv))
	case int64:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(rv), nil
	case int32:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(rv), nil
	case int16:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(rv), nil
	case int8:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(rv), nil
	case uint:
		return rv, nil
	case uint64:
		return uint(rv), nil
	case uint32:
		return uint(rv), nil
	case uint16:
		return uint(rv), nil
	case uint8:
		return uint(rv), nil
	case float64:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(rv), nil
	case float32:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(rv), nil
	case bool:
		if rv {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint", i, i)
	}
}

func ToUint64E(i interface{}) (uint64, error) {
	i = reflects.Indirect(i)

	intv, ok := toInt(i)
	if ok {
		if intv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(intv), nil
	}

	switch rv := i.(type) {
	case string:
		v, err := strconv.ParseUint(numbers.TrimZeroDecimal(rv), 0, 0)
		if err == nil {
			if v < 0 {
				return 0, errNegativeNotAllowed
			}
			return v, nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint64", i, i)
	case json.Number:
		return ToUint64E(string(rv))
	case int64:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(rv), nil
	case int32:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(rv), nil
	case int16:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(rv), nil
	case int8:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(rv), nil
	case uint:
		return uint64(rv), nil
	case uint64:
		return rv, nil
	case uint32:
		return uint64(rv), nil
	case uint16:
		return uint64(rv), nil
	case uint8:
		return uint64(rv), nil
	case float32:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(rv), nil
	case float64:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(rv), nil
	case bool:
		if rv {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint64", i, i)
	}
}

func ToUint32E(i interface{}) (uint32, error) {
	i = reflects.Indirect(i)

	intv, ok := toInt(i)
	if ok {
		if intv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(intv), nil
	}

	switch rv := i.(type) {
	case string:
		v, err := strconv.ParseInt(numbers.TrimZeroDecimal(rv), 0, 0)
		if err == nil {
			if v < 0 {
				return 0, errNegativeNotAllowed
			}
			return uint32(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint32", i, i)
	case json.Number:
		return ToUint32E(string(rv))
	case int64:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(rv), nil
	case int32:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(rv), nil
	case int16:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(rv), nil
	case int8:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(rv), nil
	case uint:
		return uint32(rv), nil
	case uint64:
		return uint32(rv), nil
	case uint32:
		return rv, nil
	case uint16:
		return uint32(rv), nil
	case uint8:
		return uint32(rv), nil
	case float64:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(rv), nil
	case float32:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(rv), nil
	case bool:
		if rv {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint32", i, i)
	}
}

func ToUint16E(i interface{}) (uint16, error) {
	i = reflects.Indirect(i)

	intv, ok := toInt(i)
	if ok {
		if intv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(intv), nil
	}

	switch rv := i.(type) {
	case string:
		v, err := strconv.ParseInt(numbers.TrimZeroDecimal(rv), 0, 0)
		if err == nil {
			if v < 0 {
				return 0, errNegativeNotAllowed
			}
			return uint16(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint16", i, i)
	case json.Number:
		return ToUint16E(string(rv))
	case int64:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(rv), nil
	case int32:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(rv), nil
	case int16:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(rv), nil
	case int8:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(rv), nil
	case uint:
		return uint16(rv), nil
	case uint64:
		return uint16(rv), nil
	case uint32:
		return uint16(rv), nil
	case uint16:
		return rv, nil
	case uint8:
		return uint16(rv), nil
	case float64:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(rv), nil
	case float32:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(rv), nil
	case bool:
		if rv {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint16", i, i)
	}
}

func ToUint8E(i interface{}) (uint8, error) {
	i = reflects.Indirect(i)

	intv, ok := toInt(i)
	if ok {
		if intv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(intv), nil
	}

	switch rv := i.(type) {
	case string:
		v, err := strconv.ParseInt(numbers.TrimZeroDecimal(rv), 0, 0)
		if err == nil {
			if v < 0 {
				return 0, errNegativeNotAllowed
			}
			return uint8(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint8", i, i)
	case json.Number:
		return ToUint8E(string(rv))
	case int64:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(rv), nil
	case int32:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(rv), nil
	case int16:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(rv), nil
	case int8:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(rv), nil
	case uint:
		return uint8(rv), nil
	case uint64:
		return uint8(rv), nil
	case uint32:
		return uint8(rv), nil
	case uint16:
		return uint8(rv), nil
	case uint8:
		return rv, nil
	case float64:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(rv), nil
	case float32:
		if rv < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(rv), nil
	case bool:
		if rv {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint8", i, i)
	}
}

func ToStringE(i interface{}) (string, error) {
	i = reflects.IndirectToStringerOrError(i)

	switch rv := i.(type) {
	case string:
		return rv, nil
	case bool:
		return strconv.FormatBool(rv), nil
	case float64:
		return strconv.FormatFloat(rv, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(rv), 'f', -1, 32), nil
	case int:
		return strconv.FormatInt(int64(rv), 10), nil
	case int64:
		return strconv.FormatInt(rv, 10), nil
	case int32:
		return strconv.FormatInt(int64(rv), 10), nil
	case int16:
		return strconv.FormatInt(int64(rv), 10), nil
	case int8:
		return strconv.FormatInt(int64(rv), 10), nil
	case uint:
		return strconv.FormatUint(uint64(rv), 10), nil
	case uint64:
		return strconv.FormatUint(uint64(rv), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(rv), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(rv), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(rv), 10), nil
	case json.Number:
		return string(rv), nil
	case []byte:
		return string(rv), nil
	case template.HTML:
		return string(rv), nil
	case template.URL:
		return string(rv), nil
	case template.JS:
		return string(rv), nil
	case template.CSS:
		return string(rv), nil
	case template.HTMLAttr:
		return string(rv), nil
	case nil:
		return "", nil
	case fmt.Stringer:
		return rv.String(), nil
	case error:
		return rv.Error(), nil
	}

	iReflectVal := reflect.ValueOf(i)
	switch iReflectVal.Kind() {
	case reflect.String:
		return iReflectVal.String(), nil
	default:
		return "", fmt.Errorf("unable to cast %#v of type %T to string", i, i)
	}
}

func ToStringMapStringE(i interface{}) (map[string]string, error) {
	m := map[string]string{}

	switch rv := i.(type) {
	case map[string]string:
		return rv, nil
	case map[string]interface{}:
		for k, val := range rv {
			m[ToString(k)] = ToString(val)
		}
		return m, nil
	case map[interface{}]string:
		for k, val := range rv {
			m[ToString(k)] = ToString(val)
		}
		return m, nil
	case map[interface{}]interface{}:
		for k, val := range rv {
			m[ToString(k)] = ToString(val)
		}
		return m, nil
	case string:
		err := jsons.JsonStringToObject(rv, &m)
		return m, err
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]string", i, i)
	}
}

func ToStringMapStringSliceE(i interface{}) (map[string][]string, error) {
	m := map[string][]string{}

	switch rv := i.(type) {
	case map[string][]string:
		return rv, nil
	case map[string][]interface{}:
		for k, val := range rv {
			m[ToString(k)] = ToStringSlice(val)
		}
		return m, nil
	case map[string]string:
		for k, val := range rv {
			m[ToString(k)] = []string{val}
		}
	case map[string]interface{}:
		for k, val := range rv {
			switch vt := val.(type) {
			case []interface{}:
				m[ToString(k)] = ToStringSlice(vt)
			case []string:
				m[ToString(k)] = vt
			default:
				m[ToString(k)] = []string{ToString(val)}
			}
		}
		return m, nil
	case map[interface{}][]string:
		for k, val := range rv {
			m[ToString(k)] = ToStringSlice(val)
		}
		return m, nil
	case map[interface{}]string:
		for k, val := range rv {
			m[ToString(k)] = ToStringSlice(val)
		}
		return m, nil
	case map[interface{}][]interface{}:
		for k, val := range rv {
			m[ToString(k)] = ToStringSlice(val)
		}
		return m, nil
	case map[interface{}]interface{}:
		for k, val := range rv {
			key, err := ToStringE(k)
			if err != nil {
				return m, fmt.Errorf("unable to cast %#v of type %T to map[string][]string", i, i)
			}
			value, err := ToStringSliceE(val)
			if err != nil {
				return m, fmt.Errorf("unable to cast %#v of type %T to map[string][]string", i, i)
			}
			m[key] = value
		}
	case string:
		err := jsons.JsonStringToObject(rv, &m)
		return m, err
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string][]string", i, i)
	}
	return m, nil
}

func ToStringMapBoolE(i interface{}) (map[string]bool, error) {
	m := map[string]bool{}

	switch rv := i.(type) {
	case map[interface{}]interface{}:
		for k, val := range rv {
			m[ToString(k)] = ToBool(val)
		}
		return m, nil
	case map[string]interface{}:
		for k, val := range rv {
			m[ToString(k)] = ToBool(val)
		}
		return m, nil
	case map[string]bool:
		return rv, nil
	case string:
		err := jsons.JsonStringToObject(rv, &m)
		return m, err
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]bool", i, i)
	}
}

func ToStringMapE(i interface{}) (map[string]interface{}, error) {
	m := map[string]interface{}{}

	switch rv := i.(type) {
	case map[interface{}]interface{}:
		for k, val := range rv {
			m[ToString(k)] = val
		}
		return m, nil
	case map[string]interface{}:
		return rv, nil
	case string:
		err := jsons.JsonStringToObject(rv, &m)
		return m, err
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]interface{}", i, i)
	}
}

func ToStringMapIntE(i interface{}) (map[string]int, error) {
	m := map[string]int{}
	if i == nil {
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]int", i, i)
	}

	switch rv := i.(type) {
	case map[interface{}]interface{}:
		for k, val := range rv {
			m[ToString(k)] = ToInt(val)
		}
		return m, nil
	case map[string]interface{}:
		for k, val := range rv {
			m[k] = ToInt(val)
		}
		return m, nil
	case map[string]int:
		return rv, nil
	case string:
		err := jsons.JsonStringToObject(rv, &m)
		return m, err
	}

	if reflect.TypeOf(i).Kind() != reflect.Map {
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]int", i, i)
	}

	mVal := reflect.ValueOf(m)
	v := reflect.ValueOf(i)
	for _, keyVal := range v.MapKeys() {
		val, err := ToIntE(v.MapIndex(keyVal).Interface())
		if err != nil {
			return m, fmt.Errorf("unable to cast %#v of type %T to map[string]int", i, i)
		}
		mVal.SetMapIndex(keyVal, reflect.ValueOf(val))
	}
	return m, nil
}

func ToStringMapInt64E(i interface{}) (map[string]int64, error) {
	m := map[string]int64{}
	if i == nil {
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]int64", i, i)
	}

	switch rv := i.(type) {
	case map[interface{}]interface{}:
		for k, val := range rv {
			m[ToString(k)] = ToInt64(val)
		}
		return m, nil
	case map[string]interface{}:
		for k, val := range rv {
			m[k] = ToInt64(val)
		}
		return m, nil
	case map[string]int64:
		return rv, nil
	case string:
		err := jsons.JsonStringToObject(rv, &m)
		return m, err
	}

	if reflect.TypeOf(i).Kind() != reflect.Map {
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]int64", i, i)
	}
	mVal := reflect.ValueOf(m)
	v := reflect.ValueOf(i)
	for _, keyVal := range v.MapKeys() {
		val, err := ToInt64E(v.MapIndex(keyVal).Interface())
		if err != nil {
			return m, fmt.Errorf("unable to cast %#v of type %T to map[string]int64", i, i)
		}
		mVal.SetMapIndex(keyVal, reflect.ValueOf(val))
	}
	return m, nil
}

func ToSliceE(i interface{}) ([]interface{}, error) {
	var s []interface{}

	switch rv := i.(type) {
	case []interface{}:
		return append(s, rv...), nil
	case []map[string]interface{}:
		for _, u := range rv {
			s = append(s, u)
		}
		return s, nil
	default:
		return s, fmt.Errorf("unable to cast %#v of type %T to []interface{}", i, i)
	}
}

func ToBoolSliceE(i interface{}) ([]bool, error) {
	if i == nil {
		return []bool{}, fmt.Errorf("unable to cast %#v of type %T to []bool", i, i)
	}

	switch rv := i.(type) {
	case []bool:
		return rv, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]bool, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := ToBoolE(s.Index(j).Interface())
			if err != nil {
				return []bool{}, fmt.Errorf("unable to cast %#v of type %T to []bool", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []bool{}, fmt.Errorf("unable to cast %#v of type %T to []bool", i, i)
	}
}

func ToStringSliceE(i interface{}) ([]string, error) {
	var a []string

	switch rv := i.(type) {
	case []interface{}:
		for _, u := range rv {
			a = append(a, ToString(u))
		}
		return a, nil
	case []string:
		return rv, nil
	case []int8:
		for _, u := range rv {
			a = append(a, ToString(u))
		}
		return a, nil
	case []int:
		for _, u := range rv {
			a = append(a, ToString(u))
		}
		return a, nil
	case []int32:
		for _, u := range rv {
			a = append(a, ToString(u))
		}
		return a, nil
	case []int64:
		for _, u := range rv {
			a = append(a, ToString(u))
		}
		return a, nil
	case []float32:
		for _, u := range rv {
			a = append(a, ToString(u))
		}
		return a, nil
	case []float64:
		for _, u := range rv {
			a = append(a, ToString(u))
		}
		return a, nil
	case string:
		return strings.Fields(rv), nil
	case []error:
		for _, err := range i.([]error) {
			a = append(a, err.Error())
		}
		return a, nil
	case interface{}:
		str, err := ToStringE(rv)
		if err != nil {
			return a, fmt.Errorf("unable to cast %#v of type %T to []string", i, i)
		}
		return []string{str}, nil
	default:
		return a, fmt.Errorf("unable to cast %#v of type %T to []string", i, i)
	}
}

func ToIntSliceE(i interface{}) ([]int, error) {
	if i == nil {
		return []int{}, fmt.Errorf("unable to cast %#v of type %T to []int", i, i)
	}

	switch rv := i.(type) {
	case []int:
		return rv, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]int, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := ToIntE(s.Index(j).Interface())
			if err != nil {
				return []int{}, fmt.Errorf("unable to cast %#v of type %T to []int", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []int{}, fmt.Errorf("unable to cast %#v of type %T to []int", i, i)
	}
}

func ToDurationSliceE(i interface{}) ([]time.Duration, error) {
	if i == nil {
		return []time.Duration{}, fmt.Errorf("unable to cast %#v of type %T to []time.Duration", i, i)
	}

	switch rv := i.(type) {
	case []time.Duration:
		return rv, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]time.Duration, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := ToDurationE(s.Index(j).Interface())
			if err != nil {
				return []time.Duration{}, fmt.Errorf("unable to cast %#v of type %T to []time.Duration", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []time.Duration{}, fmt.Errorf("unable to cast %#v of type %T to []time.Duration", i, i)
	}
}

func toInt(v interface{}) (int64, bool) {
	val := reflect.ValueOf(v)
	if val.CanInt() {
		return val.Int(), true
	}
	switch v := v.(type) {
	case int:
		return int64(v), true
	case time.Weekday:
		return int64(v), true
	case time.Month:
		return int64(v), true
	default:
		return 0, false
	}
}
