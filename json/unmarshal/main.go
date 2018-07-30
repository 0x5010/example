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

func main() {
	input := `
{
    "type": "sound",
    "msg": {
        "description": "dynamite",
        "authority": "the Bruce Dickinson"
    }
}
`
	var env Envelope
	if err := json.Unmarshal([]byte(input), &env); err != nil {
		log.Fatal(err)
	}
	// for the love of Gopher DO NOT DO THIS
	desc := env.Msg.(map[string]interface{})["description"].(string)
	fmt.Println(desc)
	// dynamite
}
