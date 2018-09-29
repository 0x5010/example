package main

import (
	"io"
	"log"
	"net"
	"sync"

	"github.com/hashicorp/yamux"
)

const LoopCount = 1000

func main() {
	log.Println("Starting yamux demo")

	localAddr := "127.0.0.1:4444"
	wg := &sync.WaitGroup{}
	go server(localAddr, wg)

	if err := yclient(localAddr); err != nil {
		log.Println(err)
	}
	wg.Wait()
}

func yclient(serverAddr string) error {
	// Get a TCP connection
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return err
	}

	// Setup client side of yamux
	log.Println("creating client session")
	session, err := yamux.Client(conn, nil)
	if err != nil {
		return err
	}

	for i := 0; i < 100; i++ {
		// Open a new stream
		log.Println("opening stream")
		stream, err := session.Open()
		if err != nil {
			log.Print(err)
			return err
		}
		go func() {
			// Stream implements net.Conn
			for i := 0; i < LoopCount; i++ {
				_, err = stream.Write([]byte("ping"))
			}
		}()
	}

	// Open a new stream
	log.Println("opening stream2")
	stream2, err := session.Open()
	if err != nil {
		return err
	}
	// Stream implements net.Conn
	_, err = stream2.Write([]byte("ping2"))

	// Accept a stream
	log.Println("accepting stream3")
	stream3, err := session.Accept()
	if err != nil {
		return err
	}
	// Listen for a message
	buf2 := make([]byte, 5)
	_, err = stream3.Read(buf2)
	log.Printf("buf3 = %+v\n", string(buf2))

	return err
}

func server(localAddr string, wg *sync.WaitGroup) error {
	// Accept a TCP connection
	listener, err := net.Listen("tcp", localAddr)

	conn, err := listener.Accept()
	if err != nil {
		return err
	}

	wg.Add(1)
	// Setup server side of yamux
	log.Println("creating server session")
	session, err := yamux.Server(conn, nil)
	if err != nil {
		return err
	}

	for i := 0; i < 100; i++ {
		// Accept a stream
		log.Println("accepting stream")
		stream, err := session.Accept()
		if err != nil {
			log.Print(err)
			return err
		}
		go func(i int) {
			wg.Add(1)
			buf1 := make([]byte, 4)
			// Stream implements net.Conn
			for i := 0; i < LoopCount; i++ {
				// Listen for a message
				n, err := stream.Read(buf1)
				if err != nil || n != len(buf1) || string(buf1) != "ping" {
					log.Printf("read error: %s", err)
				}
			}
			wg.Done()
			log.Printf("buf#%d read done", i)
		}(i)
		wg.Add(1)
		tmp := make([]byte, 4)
		n, err := stream.Read(tmp)
		if err == io.EOF {
			log.Print("eof", n)
		}
		wg.Done()
	}

	// Accept a stream
	log.Println("accepting stream2")
	stream2, err := session.Accept()
	if err != nil {
		return err
	}
	// Listen for a message
	buf2 := make([]byte, 5)
	_, err = stream2.Read(buf2)
	log.Printf("buf2 = %+v\n", string(buf2))

	// Open a new stream
	log.Println("opening stream3")
	stream3, err := session.Open()
	if err != nil {
		return err
	}
	// Stream implements net.Conn
	_, err = stream3.Write([]byte("pong3"))

	wg.Done()
	return err
}
