package runtime

import "reflect"

func IsZero(x interface{}) bool {
	if x == nil {
		return true
	}
	value := reflect.ValueOf(x)
	switch value.Kind() {
	case reflect.Map:
		return value.Len() == 0
	case reflect.Slice:
		return value.Len() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
