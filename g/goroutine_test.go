package g

import (
	"fmt"
	"testing"
	"time"
)

func Test_pool_JJ(t *testing.T) {
	poolJJ(100, 100)
	poolJJ(50, 100)
}

func poolJJ(workerCount, resLength int) {
	t1 := time.Now()

	defer func() {
		fmt.Printf("the poolJJ workerCount:%d,resLength:%d,use time %v  \n", workerCount, resLength, time.Since(t1))
	}()

	p := NewPool(workerCount, 1)

	length := resLength
	res := make([]int, length)

	for index := range res {
		index := index
		p.JJ(func() error {
			time.Sleep(time.Second)
			res[index] = index
			return nil
		})
	}
	err := p.Wait()
	if err != nil {
		panic(err)
	}
	for index, item := range res {
		if index != item {
			panic("index != item ")
			return
		}
	}
}
