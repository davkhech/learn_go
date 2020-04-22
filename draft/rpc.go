package main

import (
	"fmt"
	"net"
	"net/rpc"
	"strconv"
)

type Server struct {}
func (s *Server) Negate(i int, reply *int) error {
	*reply = -i
	return nil
}

func serverRPC() {
	_ = rpc.Register(new(Server))
	ln, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		c, err := ln.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(c)
	}
}

func clientRPC(num int) {
	c, err := rpc.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return
	}

	var result int
	err = c.Call("Server.Negate", num, &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Server.Negate returned ", result)
}

func main() {
	go serverRPC()
	go clientRPC(999)

	var msg string
	fmt.Scanln(&msg)
	neg, err := strconv.Atoi(msg)
	if err != nil {
		fmt.Println("Provided string is not an integer")
	}
	go clientRPC(neg)

	fmt.Scanln(&msg)
}
