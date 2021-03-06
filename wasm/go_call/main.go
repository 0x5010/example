package main

import (
	"fmt"
	"syscall/js"
)

var done = make(chan struct{})

func main() {
	callback := js.NewCallback(printMessage)
	defer callback.Close()

	setPrintMessage := js.Global().Get("setPrintMessage")
	setPrintMessage.Invoke(callback)
	<-done
}

func printMessage(args []js.Value) {
	message := args[0].String()
	fmt.Println(message)
	done <- struct{}{}
}
