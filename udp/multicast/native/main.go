package main

import (
	"fmt"
	"net"

	"golang.org/x/net/ipv4"
)

func main() {
	en0, err := net.InterfaceByName("en0")
	if err != nil {
		fmt.Println(err.Error())
	}
	group := net.IPv4(224, 0, 0, 250)

	c, err := net.ListenPacket("udp4", "0.0.0.0:1024")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer c.Close()

	p := ipv4.NewPacketConn(c)
	if err := p.JoinGroup(en0, &net.UDPAddr{IP: group}); err != nil {
		fmt.Println(err.Error())
	}

	if err := p.SetControlMessage(ipv4.FlagDst, true); err != nil {
		fmt.Println(err.Error())
	}

	b := make([]byte, 1500)
	for {
		n, cm, src, err := p.ReadFrom(b)
		if err != nil {
			fmt.Println(err.Error())
		}
		if cm.Dst.IsMulticast() {
			if cm.Dst.Equal(group) {
				fmt.Printf("received: %s from <%s>\n", b[:n], src)
				n, err = p.WriteTo([]byte("world"), cm, src)
				if err != nil {
					fmt.Println(err.Error())
				}
			} else {
				fmt.Println("Unknown group")
				continue
			}
		}
	}

}
