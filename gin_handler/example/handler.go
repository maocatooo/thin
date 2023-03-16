package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/maocatooo/thin/gin_handler"
)

type Req struct {
	Name string `json:"name"`
}

type Query struct {
	Name string `form:"name"`
}

type Resp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

/*
GET http://127.0.0.1:8080/ping?name=123
*/
func Ping(ctx *gin.Context, req *Query, rsp *Resp) error {
	fmt.Printf("Ping req: %+v \n", *req)
	if req.Name == "123" {
		return fmt.Errorf("err 123")
	}
	rsp.Code = 200
	rsp.Message = req.Name
	return nil
}

type A struct {
	a string
}

/*
POST http://127.0.0.1:8080/pong
{
	"name":"456"
}
*/
func (a A) Pong(ctx *gin.Context, req *Req, rsp *Resp) error {
	fmt.Printf("Pong req: %+v \n", *req)
	if req.Name == "123" {
		return fmt.Errorf("123")
	}
	rsp.Code = 200
	rsp.Message = req.Name
	return nil
}

func main() {
	r := gin.Default()
	r.GET("/ping", gin_handler.Query(Ping))
	r.POST("/pong", gin_handler.JSON(A{a: "a123"}.Pong))
	_ = r.Run()
}
