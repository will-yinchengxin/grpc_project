package util

import (
	"fmt"
	"net"
)

func GetAddr() string {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	return parsePortToAddr(port)
}

func parsePortToAddr(port int) string {
	return fmt.Sprintf("0.0.0.0:%d", 56043)
}
