package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/hashicorp/yamux"
)

const localAddr = "127.0.0.1:4444"

func main() {
	go server()
	// wg := sync.WaitGroup{}
	// wg.Add(1)
	client()
	select {}
}

func client() {
	// Get a TCP connection
	conn, err := net.Dial("tcp", localAddr)
	if err != nil {
		panic(err)
	}

	// Setup client side of yamux
	session, err := yamux.Client(conn, nil)
	if err != nil {
		panic(err)
	}

	// Open a new stream
	for i := 0; i < 10; i++ {
		go func() {
			stream, err := session.Open()
			if err != nil {
				panic(err)
				// Stream implements net.Conn
			}
			for j := 0; j < 3; j++ {
				for i := 0; i < 3; i++ {
					stream.Write([]byte("ping"))
				}
				time.Sleep(1 * time.Second)
			}

			stream.Close()
		}()
	}
}

func server() {
	listener, err := net.Listen("tcp", localAddr)
	// Accept a TCP connection
	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}

	// Setup server side of yamux
	session, err := yamux.Server(conn, nil)
	if err != nil {
		panic(err)
	}
	n := 0
	// Accept a stream
	for {
		stream, err := session.Accept()
		n++
		if err != nil {
			panic(err)
		}

		// Listen for a message
		go func(stream net.Conn, nn int) {
			buf := make([]byte, 1024)
			for {
				n, err := stream.Read(buf)
				if err == io.EOF {

					break
				}
				fmt.Println(n, err, nn)
			}
			fmt.Println("end", nn)
		}(stream, n)
	}

}
