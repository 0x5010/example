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

		// omit keys
		OmitAuthority string `json:"authority,omitempty"`

		// add keys
		Authority string `json:"author"`
	}{
		Sound:     s,
		Authority: s.Authority,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", buf)
}
