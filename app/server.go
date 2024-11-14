package main

import (
	"errors"
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/store"
	"io"
	"net"
	"os"
	"strings"
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

	fmt.Println("Listening on port 6379")

	str := store.NewStore()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection")
			os.Exit(1)
		}

		go func() {
			err := readMultipleCommands(conn, str)
			if err != nil {
				return
			}
		}()
	}
}

func readMultipleCommands(conn net.Conn, str *store.Store) error {

	buf := make([]byte, 1024)

	for {
		_, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Failed to read from connection")
				return err
			}
			break
		}
		toWrite, err := redisProtocolParser(buf, str)
		if err != nil {
			fmt.Println("Failed to parse command")
			os.Exit(1)
		}
		_, err = conn.Write([]byte(toWrite))
		if err != nil {
			fmt.Println("Failed to write to connection")
		}

	}
	return nil

}

func redisProtocolParser(buf []byte, store *store.Store) (string, error) {
	str := string(buf)
	arrStr := strings.Split(str, "\r\n")

	command := arrStr[2]
	command = strings.ToUpper(command)
	toWrite := ""
	switch command {
	case "PING":
		toWrite, _ = handlePing()
	case "ECHO":
		toWrite, _ = handleEcho(arrStr[4])
	case "SET":
		toWrite = handleSet(arrStr[3:], store)
	case "GET":
		toWrite, _ = handleGet(arrStr[4], store)
	default:
		return "", errors.New("unknown command")
	}

	return toWrite, nil

}
