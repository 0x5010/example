package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func simple() {
	reader := strings.NewReader("Clear is better than clever")
	p := make([]byte, 4)

	for {
		n, err := reader.Read(p)
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF:", n)
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(n, string(p[:n]))
	}
}

type alphaReader struct {
	src string
	cur int
}

func newAlphaReader(src string) *alphaReader {
	return &alphaReader{src: src}
}

func alpha(r byte) byte {
	if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
		return r
	}
	return 0
}

func (a *alphaReader) Read(p []byte) (int, error) {
	if a.cur >= len(a.src) {
		return 0, io.EOF
	}
	x := len(a.src) - a.cur
	n, bound := 0, 0
	if x >= len(p) {
		bound = len(p)
	} else if x < len(p) {
		bound = x
	}

	buf := make([]byte, bound)
	for n < bound {
		if char := alpha(a.src[a.cur]); char != 0 {
			buf[n] = char
		}
		n++
		a.cur++
	}
	copy(p, buf)
	return n, nil
}

func customReader() {
	reader := newAlphaReader("Hello! It's 9am, where is the sun?")
	p := make([]byte, 4)
	for {
		n, err := reader.Read(p)
		if err == io.EOF {
			break
		}
		fmt.Print(string(p[:n]))
	}
	fmt.Println()
}

type alphaReader2 struct {
	reader io.Reader
}

func newAlphaReader2(reader io.Reader) *alphaReader2 {
	return &alphaReader2{reader: reader}
}

func (a *alphaReader2) Read(p []byte) (int, error) {
	n, err := a.reader.Read(p)
	if err != nil {
		return n, err
	}
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		if char := alpha(p[i]); char != 0 {
			buf[i] = char
		}
	}
	copy(p, buf)
	return n, nil
}

func group() {
	reader := newAlphaReader2(strings.NewReader("Hello! It's 9am, where is the sun?"))
	p := make([]byte, 4)
	for {
		n, err := reader.Read(p)
		if err == io.EOF {
			break
		}
		fmt.Print(string(p[:n]))
	}
	fmt.Println()
}

func fileGroup() {
	file, err := os.Open("./test")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	reader := newAlphaReader2(file)
	p := make([]byte, 4)
	for {
		n, err := reader.Read(p)
		if err == io.EOF {
			break
		}
		fmt.Print(string(p[:n]))
	}
	fmt.Println()
}

func main() {
	simple()
	customReader()
	group()
	fileGroup()
}
