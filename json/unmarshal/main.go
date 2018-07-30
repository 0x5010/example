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
	input := `
{
    "type": "sound",
    "msg": {
        "description": "dynamite",
        "authority": "the Bruce Dickinson"
    }
}
`
	{
		var env Envelope
		if err := json.Unmarshal([]byte(input), &env); err != nil {
			log.Fatal(err)
		}

		desc := env.Msg.(map[string]interface{})["description"].(string)
		fmt.Println(desc)
	}

	{
		var msg json.RawMessage
		env := Envelope{
			Msg: &msg,
		}
		if err := json.Unmarshal([]byte(input), &env); err != nil {
			log.Fatal(err)
		}
		switch env.Type {
		case "sound":
			var s Sound
			if err := json.Unmarshal(msg, &s); err != nil {
				log.Fatal(err)
			}
			desc := s.Description
			fmt.Println(desc)
		default:
			log.Fatalf("unknown message type: %q", env.Type)
		}
	}

}
