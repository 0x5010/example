package ginkgo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/0x5010/example/ginkgo/books"
)

var _ = Describe("Book", func() {
	var (
		longBook  Book
		shortBook Book
	)

	BeforeEach(func() {
		longBook = Book{
			Title:  "Les Miserables",
			Author: "Victor Hugo",
			Pages:  1488,
		}

		shortBook = Book{
			Title:  "Fox In Socks",
			Author: "Dr. Seuss",
			Pages:  24,
		}
	})

	Describe("Categorizing book length", func() {
		Context("With more than 300 pages", func() {
			It("should be a novel", func() {
				Expect(longBook.CategoryByLength()).To(Equal("NOVEL"))
			})
		})

		Context("With fewer than 300 pages", func() {
			It("should be a short story", func() {
				Expect(shortBook.CategoryByLength()).To(Equal("SHORT STORY"))
			})
		})

		It("panic in a goroutine", func(done Done) {
			go func() {
				defer GinkgoRecover()
				Ω(DoSomething()).Should(BeTrue())
				close(done)
			}()
		})

	})
})

var _ = Describe("Book", func() {
	var (
		book Book
		err  error
	)

	BeforeEach(func() {
		book, err = NewBookFromJSON(`{
            "title":"Les Miserables",
            "author":"Victor Hugo",
            "pages":1488
        }`)
	})

	It("can be loaded from JSON", func() {
		Expect(book.Title).To(Equal("Les Miserables"))
		Expect(book.Author).To(Equal("Victor Hugo"))
		Expect(book.Pages).To(Equal(1488))
	})

	It("can extract the author's last name", func() {
		Expect(book.AuthorLastName()).To(Equal("Hugo"))
	})
})
