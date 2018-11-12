package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/observable"
	"github.com/reactivex/rxgo/observer"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	lastPos := 0
	pchan := make(chan interface{})

	go func() {
		for {
			select {
			case <-time.After(2 * time.Second):
				log.Printf("User has stopped @%dpx\n", lastPos)
			case pos := <-pchan:
				log.Printf("User @%dpx\n", pos.(int))
			}
		}
	}()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		source := observable.Just(p)

		onNext := handlers.NextFunc(func(item interface{}) {
			if item, ok := item.([]byte); ok {
				offset, _ := strconv.Atoi(string(item))
				pchan <- offset
				lastPos = offset
			}
		})

		_ = source.Subscribe(observer.New(onNext))
	}
}

func main() {
	http.HandleFunc("/scroll", indexHandler)
	http.Handle("/", http.FileServer(http.Dir(".")))

	log.Println("Server listening on port 4000")
	log.Fatal(http.ListenAndServe(":4000", nil))
}
