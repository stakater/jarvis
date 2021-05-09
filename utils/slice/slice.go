package slice

import "reflect"

func Contains(inputSlice interface{}, element interface{}) bool {

	v := reflect.ValueOf(inputSlice)

	if v.Kind() == reflect.Slice {
		length := v.Len()
		for i := 0; i < length; i++ {
			ele := v.Index(i).Interface()
			if ele == element {
				return true
			}
		}
	}
	return false
}
