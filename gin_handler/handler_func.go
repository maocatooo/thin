package gin_handler

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

var Post = HandlerFunc

func HandlerFunc(hd interface{}) gin.HandlerFunc {
	handler, err := toHandler(hd)
	if err != nil {
		panic(err)
	}
	return func(ctx *gin.Context) {
		var (
			req = reflect.New(handler.reqType.Elem())
			rsp = reflect.New(handler.rspType.Elem())
		)

		if err := ctx.BindJSON(req.Interface()); err != nil {
			_ = ctx.AbortWithError(400, err)
			return
		}
		if err := handler.call(ctx, req, rsp); err != nil {
			_ = ctx.AbortWithError(400, err)
			return
		}

		ctx.AbortWithStatusJSON(200, rsp.Interface())
	}
}
