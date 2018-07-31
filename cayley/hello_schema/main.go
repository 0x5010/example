package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"

	"github.com/cayleygraph/cayley"

	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/kv/leveldb"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/schema"
	"github.com/cayleygraph/cayley/voc"
)

type Person struct {
	rdfType struct{} `quad:"@type > ex:Person"`
	ID      quad.IRI `json:"@id"`
	Name    string   `json:"ex:name"`
	Age     int      `quad:"ex:age"`
}

type Coords struct {
	Lat float64 `json:"ex:lat"`
	Lng float64 `json:"ex:lng"`
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	voc.RegisterPrefix("ex:", "http://example.org/")
	schema.RegisterType(quad.IRI("ex:Coords"), Coords{})
	sch := schema.NewConfig()
	sch.GenerateID = func(_ interface{}) quad.Value {
		return quad.BNode(fmt.Sprintf("node%d", rand.Intn(1000)))
	}

	tmpdir, err := ioutil.TempDir("", "example")
	checkErr(err)
	defer os.RemoveAll(tmpdir)

	err = graph.InitQuadStore("leveldb", tmpdir, nil)
	checkErr(err)

	store, err := cayley.NewGraph("leveldb", tmpdir, nil)
	checkErr(err)
	defer store.Close()
	qw := graph.NewWriter(store)

	bob := Person{
		ID:   quad.IRI("ex:bob").Full().Short(),
		Name: "Bob",
		Age:  32,
	}
	fmt.Printf("saving: %+v\n", bob)
	id, err := sch.WriteAsQuads(qw, bob)
	checkErr(err)
	err = qw.Close()
	checkErr(err)

	fmt.Println("id for object:", id, "=", bob.ID)

	var someone Person
	err = sch.LoadTo(nil, store, &someone, id)
	checkErr(err)
	fmt.Printf("loaded: %+v\n", someone)

	var people []Person
	err = sch.LoadTo(nil, store, &people)
	checkErr(err)
	fmt.Printf("people: %+v\n", people)

	fmt.Println()

	coords := []Coords{
		{Lat: 12.3, Lng: 34.5},
		{Lat: 39.7, Lng: 8.41},
	}
	qw = graph.NewWriter(store)
	for _, c := range coords {
		id, err = sch.WriteAsQuads(qw, c)
		checkErr(err)
		fmt.Println("generated id:", id)
	}
	err = qw.Close()
	checkErr(err)

	var newCoords []Coords
	err = sch.LoadTo(nil, store, &newCoords)
	checkErr(err)
	fmt.Printf("coords: %+v\n", newCoords)

	fmt.Println("\nquads:")
	ctx := context.TODO()
	it := store.QuadsAllIterator()
	for it.Next(ctx) {
		fmt.Println(store.Quad(it.Result()))
	}
}
