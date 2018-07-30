package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Envelope struct {
	Type string `json:"type"`
}

type Sound struct {
	Description string `json:"description"`
	Authority   string `json:"authority"`
}

type Cowbell struct {
	More bool `json:"more"`
}

func main() {
	input := `
{
    "type": "sound",
    "description": "dynamite",
    "authority": "the Bruce Dickinson"
}
`
	var env Envelope
	buf := []byte(input)
	if err := json.Unmarshal(buf, &env); err != nil {
		log.Fatal(err)
	}
	switch env.Type {
	case "sound":
		var env Envelope
		var s Sound

		if err := json.Unmarshal(buf, &struct {
			*Envelope
			*Sound
		}{&env, &s}); err != nil {
			log.Fatal(err)
		}
		// s := struct {
		// 	*Envelope
		// 	*Sound
		// }{}
		// if err := json.Unmarshal(buf, &s); err != nil {
		// 	log.Fatal(err)
		// }
		desc := s.Description
		fmt.Println(desc)
	default:
		log.Fatalf("unknown message type: %q", env.Type)
	}
}
