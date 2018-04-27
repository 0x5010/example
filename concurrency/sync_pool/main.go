package main

import (
	"fmt"
	"sync"
)

func main() {
	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating new instance")
			return struct{}{}
		},
	}
	myPool.Put(myPool.New())
	myPool.Get()
	instance := myPool.Get()
	fmt.Println()
	myPool.Put(instance)
	myPool.Get()
}
