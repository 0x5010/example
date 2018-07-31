package main

import (
	"fmt"

	"github.com/tylertreat/BoomFilters"
)

func main() {
	sbf := boom.NewDefaultStableBloomFilter(10000, 0.01)
	fmt.Println("stable point", sbf.StablePoint())

	sbf.Add([]byte("a"))
	if sbf.Test([]byte("a")) {
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
