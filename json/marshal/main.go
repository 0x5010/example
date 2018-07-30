package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Envelope struct {
	Type string      `json:"type"`
	Msg  interface{} `json:"msg"`
}

type Sound struct {
	Description string `json:"description"`
	Authority   string `json:"authority"`
}

type Cowbell struct {
	More bool `json:"more"`
}

func main() {
	s := Envelope{
		Type: "sound",
		Msg: &Sound{
			Description: "dynamite",
			Authority:   "the Bruce Dickinson",
		},
	}
	buf, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", buf)

	c := Envelope{
		Type: "cowbell",
		Msg: &Cowbell{
			More: true,
		},
	}
	buf, err = json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", buf)
}
