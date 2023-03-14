package gin_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io"
	"reflect"
)

// 注册POST,PUT,DELETE请求获取 JsonBody
func JsonBody(hd interface{}) gin.HandlerFunc {
	return handlerFunc(hd, binding.JSON)
}

// 注册Get请求获取 ? 之后的参数
func Query(hd interface{}) gin.HandlerFunc {
	return handlerFunc(hd, binding.Query)
}

func HandlerFunc(hd interface{}, bb binding.Binding) gin.HandlerFunc {
	return handlerFunc(hd, bb)
}

func handlerFunc(hd interface{}, bb binding.Binding) gin.HandlerFunc {
	handler, err := toHandler(hd)
	if err != nil {
		panic(err)
	}
	return func(ctx *gin.Context) {
		var (
			req = reflect.New(handler.reqType.Elem())
			rsp = reflect.New(handler.rspType.Elem())
		)
		if err := ctx.ShouldBindWith(req.Interface(), bb); err != nil && err != io.EOF {
			AbortWithBindErr(ctx, err)
			return
		}
		if err := handler.call(ctx, req, rsp); err != nil {
			AbortWithErr(ctx, err)
			return
		}

		ctx.AbortWithStatusJSON(StatusOKCode, rsp.Interface())
	}
}
