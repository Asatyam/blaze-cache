package main

import (
	"fmt"
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

func (app *application) handleGet(key string) (string, error) {

	value, found := app.store.Get(key)
	toWrite := fmt.Sprint("$-1\r\n")
	if found {
		length := len(value)
		toWrite = fmt.Sprintf("$%d\r\n%s\r\n", length, value)
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
