package broker

import (
	"context"
	"fmt"
	"testing"
)

type S struct {
	A string
}

func (m S) Do(_ context.Context, _ int) error {
	return nil
}

func (m S) Do2(_ context.Context, _ int) error {
	return nil
}

func F(_ context.Context, c int) error {
	return fmt.Errorf("%d", c)
}

func Test_methodObjToHandler(t *testing.T) {
	r1, _ := methodObjToHandler(F)
	if len(r1) != 1 {
		t.Error("methodObjToHandler F func length must 1")
		return
	}
	r2, _ := methodObjToHandler(S{})
	if len(r2) != 2 {
		t.Error("methodObjToHandler S struct length must 2")
		return
	}
}

func Test_methodObjToHandler_call_func_type(t *testing.T) {
	subs, _ := methodObjToHandler(F)
	for _, item := range subs {
		err := item.call(context.Background(), 13333)
		if err != nil && err.Error() != "13333" {
			t.Error(err)
		}

	}

}

func Test_methodObjToHandler_call_struct_type(t *testing.T) {
	subs, _ := methodObjToHandler(S{})
	for _, item := range subs {
		err := item.call(context.Background(), 13333)
		if err != nil {
			t.Error(err)
		}
	}

}

func Test_valueCopy(t *testing.T) {
	var r = 1
	rst := valueCopy(&r)
	r2, ok := rst.(*int)
	fmt.Println(*r2, ok)
}
