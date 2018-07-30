package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Sound struct {
	Description string `json:"description"`
	Authority   string `json:"authority"`
}

type omit *struct{}

type OmitSound struct {
	*Sound
	Authority omit `json:"authority,omitempty"`
}

func main() {
	s := &Sound{
		Description: "dynamite",
		Authority:   "the Bruce Dickinson",
	}

	buf, err := json.Marshal(OmitSound{
		Sound: s,
	})
	// buf, err := json.Marshal(struct {
	// 	*Sound
	// 	Authority bool `json:"authority,omitempty"`
	// }{
	// 	Sound: s,
	// })
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", buf)
}
