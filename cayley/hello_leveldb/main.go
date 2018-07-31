package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/kv/leveldb"
	"github.com/cayleygraph/cayley/quad"
)

func main() {
	tmpdir, err := ioutil.TempDir("", "example")
	if err != nil {
		log.Fatalln(err)
	}
	defer os.RemoveAll(tmpdir)

	err = graph.InitQuadStore("leveldb", tmpdir, nil)
	if err != nil {
		log.Fatal(err)
	}

	store, err := cayley.NewGraph("leveldb", tmpdir, nil)
	if err != nil {
		log.Fatal(err)
	}

	store.AddQuad(quad.Make("phrase of the day", "is of course", "Hello LevelDB!", "demo graph"))

	p := cayley.StartPath(store, quad.String("phrase of the day")).Out(quad.String("is of course"))

	it, _ := p.BuildIterator().Optimize()
	it, _ = store.OptimizeIterator(it)
	defer it.Close()

	ctx := context.TODO()
	for it.Next(ctx) {
		token := it.Result()
		value := store.NameOf(token)
		nativeValue := quad.NativeOf(value)
		fmt.Println(nativeValue)
	}
	if err := it.Err(); err != nil {
		log.Fatalln(err)
	}
}
