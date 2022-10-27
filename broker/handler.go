package broker

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"
)

const (
	subSig = "func(context.Context, interface{}) error"
)

type Subscriber interface {
	call(ctx context.Context, val interface{}) error
}

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

type handler struct {
	isFunc bool

	hd interface{}

	method reflect.Value

	structValue reflect.Value
	ctxType     reflect.Type
	reqType     reflect.Type
}

func (h *handler) call(ctx context.Context, val interface{}) error {
	var value []reflect.Value
	if !h.isFunc {
		value = append(value, h.structValue)
	}
	value = append(value,
		reflect.ValueOf(ctx),
		reflect.ValueOf(val),
	)
	returnValues := h.method.Call(value)
	if err := returnValues[0].Interface(); err != nil {
		return err.(error)
	}
	return nil
}

func methodObjToHandler(hd interface{}) ([]Subscriber, error) {
	if err := validate(hd); err != nil {
		return nil, err
	}

	var handlers []Subscriber
	// the func obj
	if typ := reflect.TypeOf(hd); typ.Kind() == reflect.Func {
		h := &handler{
			isFunc: true,
			method: reflect.ValueOf(hd),
		}
		h.ctxType = typ.In(0)
		h.reqType = typ.In(1)
		handlers = append(handlers, h)

	} else {
		// the struct obj
		structValue := reflect.ValueOf(hd) //

		for m := 0; m < typ.NumMethod(); m++ {
			method := typ.Method(m)
			h := &handler{
				method: method.Func,
			}
			h.structValue = structValue
			h.ctxType = method.Type.In(1)
			h.reqType = method.Type.In(2)
			handlers = append(handlers, h)
		}
	}
	return handlers, nil
}

func validate(hd interface{}) error {
	typ := reflect.TypeOf(hd)
	var argType reflect.Type

	if typ.Kind() == reflect.Func {
		name := "Func"
		switch typ.NumIn() {
		case 2:
			argType = typ.In(1)
		default:
			return fmt.Errorf("hd %v takes wrong number of args: %v required signature %s", name, typ.NumIn(), subSig)
		}
		if !isExportedOrBuiltinType(argType) {
			return fmt.Errorf("hd %v argument type not exported: %v", name, argType)
		}
		if typ.NumOut() != 1 {
			return fmt.Errorf("hd %v has wrong number of outs: %v require signature %s",
				name, typ.NumOut(), subSig)
		}
		if returnType := typ.Out(0); returnType != typeOfError {
			return fmt.Errorf("hd %v returns %v not error", name, returnType.String())
		}
	} else {
		hdlr := reflect.ValueOf(hd)
		name := reflect.Indirect(hdlr).Type().Name()

		for m := 0; m < typ.NumMethod(); m++ {
			method := typ.Method(m)

			switch method.Type.NumIn() {
			case 3:
				argType = method.Type.In(2)
			default:
				return fmt.Errorf("hd %v.%v takes wrong number of args: %v required signature %s",
					name, method.Name, method.Type.NumIn(), subSig)
			}

			if !isExportedOrBuiltinType(argType) {
				return fmt.Errorf("%v argument type not exported: %v", name, argType)
			}
			if method.Type.NumOut() != 1 {
				return fmt.Errorf(
					"hd %v.%v has wrong number of outs: %v require signature %s",
					name, method.Name, method.Type.NumOut(), subSig)
			}
			if returnType := method.Type.Out(0); returnType != typeOfError {
				return fmt.Errorf("hd %v.%v returns %v not error", name, method.Name, returnType.String())
			}
		}
	}

	return nil
}

// Is this an exported - upper case - name?
func isExported(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

func valueCopy(val interface{}) interface{} {
	if val == nil {
		return nil
	}
	pv := reflect.ValueOf(val)
	if pv.Kind() != reflect.Ptr || pv.IsNil() {
		return val
	}
	var bf bytes.Buffer
	err := gob.NewEncoder(&bf).Encode(val)
	if err != nil {
		panic(err)
	}
	srcV := reflect.New(pv.Elem().Type())
	i := srcV.Interface()
	err = gob.NewDecoder(&bf).Decode(i)
	if err != nil {
		panic(err)
	}
	return i
}
