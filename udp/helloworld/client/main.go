package main

import (
	"fmt"
	"net"
)

func main() {
	ip := net.ParseIP("127.0.0.1")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: ip, Port: 9999}

	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close()
	bs := make([]byte, 1600)
	conn.Write(bs)

	data := make([]byte, 1024)
	n, err := conn.Read(data)

	fmt.Printf("read %s from <%s>\n", data[:n], conn.RemoteAddr())
}
