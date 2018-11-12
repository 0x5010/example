package main

import (
	"fmt"

	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/observable"
	"github.com/reactivex/rxgo/observer"
)

func main() {
	score := 9

	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			score += num
		}
	})

	onDone := handlers.DoneFunc(func() {
		score *= 2
	})

	watcher := observer.New(onNext, onDone)
	sub := observable.Just(1).Subscribe(watcher)
	<-sub

	fmt.Println(score)
}
