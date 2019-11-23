package utils

import (
	"fmt"
	"log"
	"net"
)

func GetUDPConn(host string, port uint16) *net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	CheckErr(err)

	conn, err := net.DialUDP("udp", nil, addr)
	CheckErr(err)

	return conn
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Average(a int16, b int16) int16 {
	return (a / 2) + (b / 2) + ((a%2 + b%2) / 2)
}
