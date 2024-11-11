package main

import (
	"errors"
	"fmt"
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

	for {
		conn, err := l.Accept()
		fmt.Println("Accepting connection")
		if err != nil {
			fmt.Println("Failed to accept connection")
			os.Exit(1)
		}

		go func() {
			err := readMultipleCommands(conn)
			if err != nil {
				return
			}
		}()
		if err != nil {
			fmt.Printf("Failed to read command: %s\n", err)
			os.Exit(1)
		}
	}
}

func readMultipleCommands(conn net.Conn) error {

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
		toWrite, err := redisProtocolParser(buf)
		if err != nil {
			fmt.Println("Failed to parse command")
			os.Exit(1)
		}
		_, err = conn.Write([]byte(fmt.Sprintf("+%s\r\n", toWrite)))
		if err != nil {
			fmt.Println("Failed to write to connection")
		}

	}
	return nil

}

func redisProtocolParser(buf []byte) (string, error) {
	str := string(buf)
	arrStr := strings.Split(str, "\r\n")
	for _, v := range arrStr {
		fmt.Printf("%s\n", v)
	}
	//numParamsStr := arrStr[0][1:]
	//numParams, err := strconv.ParseInt(numParamsStr, 10, 0)
	//if err != nil {
	//	fmt.Println("Failed to parse number of parameters")
	//	return "", err
	//}

	command := arrStr[2]
	command = strings.ToUpper(command)
	switch command {
	case "PING":
		return "PONG", nil
	case "ECHO":
		value := arrStr[4]
		return value, nil
	}

	return "", errors.New("unknown command")

}
