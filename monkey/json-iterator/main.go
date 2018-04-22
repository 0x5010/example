package main

import (
	"encoding/json"
	"fmt"

	"github.com/bouk/monkey"
	jsoniter "github.com/json-iterator/go"
)

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	monkey.Patch(json.Marshal, func(v interface{}) ([]byte, error) {
		fmt.Println("use jsoniter marshal")
		return jsoniter.Marshal(v)
	})

	monkey.Patch(json.Unmarshal, func(data []byte, v interface{}) error {
		fmt.Println("use jsoniter unmarshal")
		return jsoniter.Unmarshal(data, v)
	})
	u1 := &user{
		ID:   "1",
		Name: "0x5010",
	}

	u2 := &user{}

	v, err := json.Marshal(u1)
	fmt.Println(string(v), err)

	err = json.Unmarshal(v, u2)
	fmt.Println(u2, err)
}
