package main

import (
	"fmt"
	"log"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/quad"
)

func main() {
	store, err := cayley.NewMemoryGraph()
	if err != nil {
		log.Fatalln(err)
	}
	store.AddQuad(quad.Make("phrase of the day", "is of course", "Hello World!", nil))

	p := cayley.StartPath(store, quad.String("phrase of the day")).Out(quad.String("is of course"))
	err = p.Iterate(nil).EachValue(nil, func(value quad.Value) {
		nativeValue := quad.NativeOf(value)
		fmt.Println(nativeValue)
	})
	if err != nil {
		log.Fatalln(err)
	}
}
