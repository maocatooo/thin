package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/maocatooo/thin/gin_handler"
)

type R1 struct {
	A1   string `json:"a1"`
	Name string `json:"name"`
}

type Q1 struct {
	A1   string `form:"a1"`
	Name string `form:"name"`
}

type R2 struct {
	A1 string `json:"a_1"`
	A2 string `json:"a_2"`
}

/*
POST http://127.0.0.1:8080/ping?a1=1a&&name=ddd

{
	"A1":"123"
}
*/

func A1(ctx *gin.Context, req *R1, rsp *R2) error {
	var query Q1
	err := ctx.BindQuery(&query)
	if err != nil {
		return err
	}
	fmt.Println(`req`, *req)
	fmt.Println(`query`, query)
	fmt.Println(`rsp`, rsp)
	return nil
}

func main() {

	r := gin.Default()
	r.Any("/ping", gin_handler.Post(A1))
	_ = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
