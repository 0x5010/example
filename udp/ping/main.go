package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func read(conn *net.UDPConn) {
	for {
		data := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		fmt.Printf("<%s> receive %s from <%s>\n", conn.LocalAddr().String(), data[:n], remoteAddr)
	}
}
func main() {
	addr1 := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9981}
	addr2 := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9982}
	go func() {
		listener1, err := net.ListenUDP("udp", addr1)
		if err != nil {
			fmt.Println(err)
			return
		}
		go read(listener1)
		for {
			listener1.WriteToUDP([]byte("ping to #2: "+addr2.String()), addr2)
			time.Sleep(5 * time.Second)
		}
	}()
	go func() {
		listener1, err := net.ListenUDP("udp", addr2)
		if err != nil {
			fmt.Println(err)
			return
		}
		go read(listener1)
		for {
			listener1.WriteToUDP([]byte("ping to #1: "+addr1.String()), addr1)
			time.Sleep(5 * time.Second)
		}

	}()
	b := make([]byte, 1)
	os.Stdin.Read(b)
}
