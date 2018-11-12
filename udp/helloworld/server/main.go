package main

import (
	"fmt"
	"net"
)

func main() {
	listerner, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9999})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Local: <%s> \n", listerner.LocalAddr().String())
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listerner.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		fmt.Printf("<%s> %s\n", remoteAddr, data[:n])
		_, err = listerner.WriteToUDP([]byte("world"), remoteAddr)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
