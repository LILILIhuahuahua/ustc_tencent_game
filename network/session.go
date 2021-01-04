package network

import (
	"net"
)

type session struct {
	conn *net.UDPConn
	remote *net.UDPAddr
}