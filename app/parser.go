package main

import (
	"fmt"
	"os"
	"strings"
)

func (app *application) handlePing() (string, error) {

	return "+PONG\r\n", nil
}

func (app *application) handleEcho(str string) (string, error) {

	toWrite := fmt.Sprintf("+%s\r\n", str)
	return toWrite, nil
}

func (app *application) handleSet(arrString []string) string {

	oldValue, err := app.store.Set(arrString)
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

func (app *application) inMemoryGet(key string) (string, error) {
	toWrite := fmt.Sprint("$-1\r\n")

	value, found := app.store.Get(key)
	if found {
		length := len(value)
		toWrite = fmt.Sprintf("$%d\r\n%s\r\n", length, value)
	}
	return toWrite, nil
}

func (app *application) handleGet(key string) (string, error) {

	toWrite := fmt.Sprint("$-1\r\n")

	path, err := app.getFilePath()
	if err != nil {
		toWrite, _ = app.inMemoryGet(key)
		return toWrite, nil
	}
	file, err := os.ReadFile(path)
	if err != nil {
		return "", ErrFileNotFound
	}

	_, err = app.verifyRDBFile(file)
	if err != nil {
		return "", err
	}
	contents := app.parseTable(file)
	for i := 2; i < len(contents); {
		if contents[i] == 0xff {
			break
		}
		fmt.Printf("i=%d ", i)
		i += 1
		keyLen := int(contents[i])
		currKey := string(contents[i+1 : i+keyLen+1])
		valueLen := int(contents[i+keyLen+1])
		value := contents[i+keyLen+2 : i+keyLen+valueLen+2]
		if currKey == key {
			toWrite = fmt.Sprintf("$%d\r\n%s\r\n", valueLen, value)
			return toWrite, nil
		}
		i = i + keyLen + 2 + valueLen
		fmt.Printf("i=%d key=%s\n ", i, key)

	}
	
	return toWrite, nil
}

func (app *application) handleConfig(arrString []string) (string, error) {

	cfgCommand := arrString[1]
	cfgCommand = strings.ToUpper(cfgCommand)
	toWrite := ""
	switch cfgCommand {
	case "GET":
		key := arrString[3]
		toWrite, _ = app.handleConfigGet(key)
	}
	return toWrite, nil

}

func (app *application) handleConfigGet(key string) (string, error) {

	value, err := app.config.Get(key)
	toWrite := fmt.Sprint("$-1\r\n")
	if err != nil {
		return toWrite, err
	}
	lenKey := len(key)
	lenValue := len(value)

	tempArr := []any{lenKey, key, lenValue, value}
	tempStr := "*2\r\n"
	for _, v := range tempArr {
		switch v.(type) {
		case string:
			tempStr += fmt.Sprint(v)
		case int:
			tempStr += fmt.Sprintf("$%d", v)

		}
		tempStr += "\r\n"
	}
	toWrite = tempStr
	return toWrite, nil
}

func (app *application) handleKeys(arrString []string) string {

	toWrite := fmt.Sprint("$-1\r\n")

	file, err := app.parseRDBFile()
	if err != nil {
		return toWrite
	}
	var keys []string
	for i := 2; i < len(file); {
		if file[i] == 0xff {
			break
		}
		fmt.Printf("i=%d ", i)
		i += 1
		keyLen := int(file[i])
		key := file[i+1 : i+keyLen+1]
		keys = append(keys, string(key))
		fmt.Printf("key=%s\n", key)
		valueLen := int(file[i+keyLen+1])
		i = i + keyLen + 2 + valueLen
		fmt.Printf("i=%d key=%s\n ", i, key)

	}
	fmt.Println(keys)
	lenKeys := len(keys)
	ans := fmt.Sprintf("*%d\r\n", lenKeys)
	for _, key := range keys {
		ans += fmt.Sprintf("$%d\r\n%s\r\n", len(key), key)
	}

	if arrString[1] == "*" {
		toWrite = ans
	}
	return toWrite

}
