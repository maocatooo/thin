package gin_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type RestResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	StatusOKCode  = http.StatusOK
	ErrCode       = http.StatusBadRequest
	ErrResultCode = http.StatusBadRequest
	BindErrCode   = http.StatusBadRequest
)

func Success(message ...string) RestResult {
	ms := ""
	if len(message) > 0 {
		ms = message[0]
	}
	return RestResult{
		Code:    http.StatusOK,
		Message: ms,
	}
}

func ErrResult(message ...string) RestResult {
	ms := ""
	if len(message) > 0 {
		ms = message[0]
	}
	return RestResult{
		Code:    ErrResultCode,
		Message: ms,
	}
}

func ErrResultWithErr(err error) RestResult {
	return ErrResult(err.Error())
}

var AbortWithErr = abortWithErr

func abortWithErr(ctx *gin.Context, err error) {
	r := RestResult{}
	r.Code = ErrCode
	r.Message = err.Error()
	ctx.AbortWithStatusJSON(ErrCode, r)
}

var AbortWithBindErr = abortWithBindErr

func abortWithBindErr(ctx *gin.Context, err error) {
	r := RestResult{}
	r.Code = BindErrCode
	r.Message = err.Error()
	ctx.AbortWithStatusJSON(BindErrCode, r)
}
