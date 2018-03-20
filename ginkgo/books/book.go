package books

import (
	"encoding/json"
	"strings"
)

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Pages  int    `json:"pages"`
}

func (b *Book) CategoryByLength() string {
	if b.Pages > 300 {
		return "NOVEL"
	}
	return "SHORT STORY"

}

func (b *Book) AuthorLastName() string {
	aa := strings.Split(b.Author, " ")
	return aa[1]
}

func DoSomething() bool {
	// fmt.Println("dosomething")
	return true
}

func NewBookFromJSON(bookJson string) (Book, error) {
	var book Book
	err := json.Unmarshal([]byte(bookJson), &book)
	return book, err
}
