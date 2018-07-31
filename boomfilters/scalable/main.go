package main

import (
	"fmt"

	"github.com/tylertreat/BoomFilters"
)

func main() {
	sbf := boom.NewDefaultScalableBloomFilter(0.01)

	if sbf.Add([]byte("a")).Test([]byte("a")) {
		fmt.Println("contains a")
	}

	if !sbf.TestAndAdd([]byte("b")) {
		fmt.Println("doesn't contain b")
	}

	if sbf.Test([]byte("b")) {
		fmt.Println("now it contains b!")
	}

	sbf.Reset()
}
