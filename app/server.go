package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var _ = net.Listen
var _ = os.Exit

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer func(l net.Listener) {
		_ = l.Close()
	}(l)

	//fmt.Println("Listening on port 6379")

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Failed to accept connection")
		os.Exit(1)
	}
	defer conn.Close()

	err = readMultipleCommand(conn)
	if err != nil {
		fmt.Println("Failed to read the inputs")
		os.Exit(1)
	}

}

func readMultipleCommand(conn net.Conn) error {

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		_, err := conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			return err
		}
	}
	return nil

}
