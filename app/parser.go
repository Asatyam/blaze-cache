package main

import (
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/store"
)

func handlePing() (string, error) {

	return "+PONG\r\n", nil
}

func handleEcho(str string) (string, error) {

	toWrite := fmt.Sprintf("+%s\r\n", str)
	return toWrite, nil
}

func handleSet(arrString []string, store *store.Store) string {

	oldValue, err := store.Set(arrString)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	toWrite := ""
	if oldValue == "" {
		toWrite = "+OK\r\n"
	} else {
		length := len(oldValue)
		toWrite = fmt.Sprintf("$%d\r\n%s\r\n", length, oldValue)
	}
	return toWrite

}

func handleGet(key string, store *store.Store) (string, error) {

	value, found := store.Get(key)
	toWrite := fmt.Sprint("$-1\r\n")
	if found {
		length := len(value)
		toWrite = fmt.Sprintf("$%d\r\n%s\r\n", length, value)
	}
	return toWrite, nil

}
