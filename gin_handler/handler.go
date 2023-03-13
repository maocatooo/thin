package gin_handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"unicode"
	"unicode/utf8"
)

const (
	subSig = "func(context.Context, interface{}, interface{}) error"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

type handler struct {
	isFunc bool

	hd interface{}

	method reflect.Value

	structValue reflect.Value
	ctxType     reflect.Type
	reqType     reflect.Type
	rspType     reflect.Type
}

func (h *handler) call(ctx *gin.Context, req, rsp reflect.Value) error {
	var value []reflect.Value
	if !h.isFunc {
		value = append(value, h.structValue)
	}
	value = append(value,
		reflect.ValueOf(ctx),
		req,
		rsp,
	)
	returnValues := h.method.Call(value)
	if err := returnValues[0].Interface(); err != nil {
		return err.(error)
	}
	return nil
}

func toHandler(hd interface{}) (*handler, error) {
	if err := validate(hd); err != nil {
		return nil, err
	}
	// the func obj
	if typ := reflect.TypeOf(hd); typ.Kind() == reflect.Func {
		h := &handler{
			isFunc: true,
			method: reflect.ValueOf(hd),
		}
		h.ctxType = typ.In(0)
		h.reqType = typ.In(1)
		h.rspType = typ.In(2)
		return h, nil
	} else {
		panic(fmt.Sprintf("hd %v takes wrong number of args: %v required signature %s", hd, typ.NumIn(), subSig))
	}
}

func validate(hd interface{}) error {
	typ := reflect.TypeOf(hd)
	var (
		argType1 reflect.Type
		argType2 reflect.Type
	)

	if typ.Kind() == reflect.Func {
		name := "Func"
		switch typ.NumIn() {
		case 3:
			argType1 = typ.In(1)
			argType2 = typ.In(2)
		default:
			return fmt.Errorf("hd %v takes wrong number of args: %v required signature %s", name, typ.NumIn(), subSig)
		}
		if !isExportedOrBuiltinType(argType1) {
			return fmt.Errorf("hd %v argument type not exported: %v", name, argType1)
		}
		if !isExportedOrBuiltinType(argType2) {
			return fmt.Errorf("hd %v argument type not exported: %v", name, argType2)
		}
		if typ.NumOut() != 1 {
			return fmt.Errorf("hd %v has wrong number of outs: %v require signature %s",
				name, typ.NumOut(), subSig)
		}
		if returnType := typ.Out(0); returnType != typeOfError {
			return fmt.Errorf("hd %v returns %v not error", name, returnType.String())
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
