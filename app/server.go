package main

import (
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
		fmt.Println("Accepting connection")
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
		if err != nil {
			fmt.Printf("Failed to read command: %s\n", err)
			os.Exit(1)
		}
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
	toWrite := ""
	switch command {
	case "PING":
		toWrite = "+PONG\r\n"
	case "ECHO":
		value := arrStr[4]
		toWrite = fmt.Sprintf("+%s\r\n", value)
	case "SET":
		value := handleSet(arrStr[3:], store)
		if value == "" {
			toWrite = "+OK\r\n"
		} else {
			length := len(value)
			toWrite = fmt.Sprintf("$%d\r\n%s\r\n", length, value)
		}
	case "GET":
		key := arrStr[4]
		value, found := handleGet(key, store)
		if !found {
			toWrite = fmt.Sprint("$-1\r\n")
		} else {
			length := len(value)
			toWrite = fmt.Sprintf("$%d\r\n%s\r\n", length, value)
		}
	}

	return toWrite, nil

}

func handleSet(arrString []string, store *store.Store) string {

	oldValue, _, err := store.Set(arrString)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return oldValue

}

func handleGet(key string, store *store.Store) (string, bool) {

	value, ok := store.Get(key)
	if !ok {
		return "", false
	}
	return value, true

}
