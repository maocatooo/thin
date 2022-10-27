# thin
go utils

Examples


`github.com/maocatooo/thin/broker`

单体项目发布订阅简单操作

```go

package main

import (
	"context"
	"fmt"
	"github.com/maocatooo/thin/broker"
)

type M struct {}

func (M) Call1(ctx context.Context, i int) error  {
	fmt.Println("M.Call1", i)
	return nil
}

func (M) Call2(ctx context.Context, i int) error  {
	fmt.Println("M.Call2", i)
	return nil
}
func (M) call3(ctx context.Context, i int) error  {
	fmt.Println("M.call3", i)
	return nil
}

func F(ctx context.Context, i int) error{
	fmt.Println("F", i)
	return nil
}

const topic = `topic.m_f`

func main()  {
	e, err := broker.NewEvent(topic, M{})
	if err != nil {
		panic(err)
	}
	err = e.Subscribe(F)
	if err != nil {
		panic(err)
	}
	err = e.Publish(context.Background(),1)
	if err != nil {
		panic(err)
	}
}


```