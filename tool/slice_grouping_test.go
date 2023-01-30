package tool

import (
	"fmt"
	"testing"
)

func Test_sliceGrouping(t *testing.T) {
	{
		s := []interface{}{"1", 2, 3}
		_, ss := sliceGrouping(s, 1, 1)
		fmt.Printf("%T, %#v \n", ss, ss) //[]interface{}{"1"}
	}
	{
		s := []interface{}{"1", 2, 3}
		_, s1 := sliceGrouping(s, 2, 2)
		fmt.Printf("%T, %#v \n", s1, s1) //[]interface{}{3}
	}
	{
		s := []interface{}{"1", 2, 3}
		_, s1 := sliceGrouping(s, 3, 2)
		fmt.Printf("%T, %#v \n", s1, s1) //[]interface{}{}
	}
}

func TestSliceGroupingWithHook(t *testing.T) {
	ms := 0
	_ = SliceGroupingWithHook([]int{1, 2, 3}, 1, func(i interface{}) error {
		s1 := i.([]int)
		ms += s1[0]
		return nil
	})
	if ms != 6 {
		t.Errorf("SliceGroupingWithHook ms must 6, it`s %d", ms)
	}
}
