package main

import (
	"fmt"
	"os"
)

func stdoutWrite() {
	proverbs := []string{
		"Channels orchestrate mutexes serialize\n",
		"Cgo is not Go\n",
		"Errors are values\n",
		"Don't panic\n",
	}

	for _, p := range proverbs {
		n, err := os.Stdout.Write([]byte(p))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if n != len(p) {
			fmt.Println("failed to write data")
			os.Exit(1)
		}
	}
}

type chanWriter struct {
	ch chan byte
}

func newChanWriter() *chanWriter {
	return &chanWriter{make(chan byte, 1024)}
}

func (w *chanWriter) Chan() <-chan byte {
	return w.ch
}

func (w *chanWriter) Write(p []byte) (int, error) {
	n := 0
	for _, b := range p {
		w.ch <- b
		n++
	}
	return n, nil
}

func (w *chanWriter) Close() error {
	close(w.ch)
	return nil
}

func customWriter() {
	writer := newChanWriter()
	go func() {
		defer writer.Close()
		writer.Write([]byte("Stream "))
		writer.Write([]byte("me!"))
	}()
	for c := range writer.Chan() {
		fmt.Printf("%c", c)
	}
	fmt.Println()
}

func main() {
	stdoutWrite()
	customWriter()
}
