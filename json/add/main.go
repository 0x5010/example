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

func main() {
	s := &Sound{
		Description: "dynamite",
		Authority:   "the Bruce Dickinson",
	}

	buf, err := json.Marshal(struct {
		*Sound
		Other string `json:"other"`
	}{
		Sound: s,
		Other: "other",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", buf)
}
