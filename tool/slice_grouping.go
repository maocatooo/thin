package tool

import (
	"reflect"
)

func SliceGroupingWithHook(s interface{}, count int, hook func(interface{}) error) error {
	for i := 1; ; i++ {
		ok, ss := sliceGrouping(s, i, count)
		if !ok {
			break
		}
		err := hook(ss)
		if err != nil {
			return err
		}

	}
	return nil
}

func SliceGrouping(s interface{}, index, count int) interface{} {
	_, ss := sliceGrouping(s, index, count)
	return ss
}

func sliceGrouping(s interface{}, index, count int) (bool, interface{}) {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Slice {
		panic("the s must type slice")
	}
	length := v.Len()
	i := (index - 1) * count
	j := index * count
	if j > length {
		j = length
	}
	if length == 0 || i >= length {
		return false, v.Slice(0, 0).Interface()
	}
	return true, v.Slice(i, j).Interface()
}
