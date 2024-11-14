package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/internal/config"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
	"io"
	"net"
	"os"
	"strings"
)

var _ = net.Listen
var _ = os.Exit

type application struct {
	store  *store.Store
	config *config.Config
}

func main() {

	var dir string
	var dbFileName string

	flag.StringVar(&dir, "dir", "./tmp", "the path to the directory where the RDB file is stored")
	flag.StringVar(&dbFileName, "dbfilename", "redis-starter.db", "the name of the RDB file")
	flag.Parse()

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
	cfg := config.NewConfig(dir, dbFileName)

	app := application{
		store:  str,
		config: cfg,
	}

	fmt.Printf("dir = %s\n dbfilename = %s\n", dir, dbFileName)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection")
			os.Exit(1)
		}

		go func() {
			err := app.parseRESP(conn)
			if err != nil {
				return
			}
		}()
	}
}

func (app *application) parseRESP(conn net.Conn) error {

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
		toWrite, err := app.parseRESPHelper(buf)
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

func (app *application) parseRESPHelper(buf []byte) (string, error) {
	str := string(buf)
	arrStr := strings.Split(str, "\r\n")

	command := arrStr[2]
	command = strings.ToUpper(command)
	toWrite := ""
	switch command {
	case "PING":
		toWrite, _ = app.handlePing()
	case "ECHO":
		toWrite, _ = app.handleEcho(arrStr[4])
	case "SET":
		toWrite = app.handleSet(arrStr[3:])
	case "GET":
		toWrite, _ = app.handleGet(arrStr[4])
	case "CONFIG":
		toWrite, _ = app.handleConfig(arrStr[3:])
	default:
		return "", errors.New("unknown command")
	}

	return toWrite, nil

}
