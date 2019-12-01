package utils

import (
	"fmt"
	"net"
)

const DefaultPort = 42170

func GetUDPConn(host string, port uint16) *net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	CheckErr(err)

	conn, err := net.DialUDP("udp", nil, addr)
	CheckErr(err)

	return conn
}
