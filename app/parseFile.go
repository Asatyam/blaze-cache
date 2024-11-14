package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

var (
	ErrNoDirectory    = errors.New("directory not specified")
	ErrNoDBFileName   = errors.New("dbfile not provided")
	ErrFileNotFound   = errors.New("file not found")
	ErrInvalidRDBFile = errors.New("invalid rdb file")
)

const (
	auxFieldCode     = 0xFA
	dbSelectorCode   = 0xFE
	resizeDBCode     = 0xFB
	expiryTimeSCode  = 0xFD
	expiryTimeMSCode = 0xFC
	eofCode          = 255
)

func (app *application) getFilePath() (string, error) {
	dir, err := app.config.Get("dir")
	if err != nil || dir == "" {
		return "", ErrNoDirectory
	}
	dbFileName, err := app.config.Get("dbfilename")
	if err != nil || dbFileName == "" {
		return "", ErrNoDBFileName
	}
	path := filepath.Join(dir, dbFileName)
	return path, nil
}

func (app *application) parseRDBFile() ([]byte, error) {

	path, err := app.getFilePath()
	if err != nil {
		return nil, err
	}
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, ErrFileNotFound
	}

	_, err = app.verifyRDBFile(file)
	if err != nil {
		return nil, err
	}
	content := app.parseTable(file)

	return content, nil
}

func (app *application) verifyRDBFile(data []byte) (int, error) {

	magic, version := string(data[:5]), string(data[5:9])
	if magic != "REDIS" {
		return 0, ErrInvalidRDBFile
	}
	versionNum, err := strconv.Atoi(version)
	if err != nil {
		return 0, ErrInvalidRDBFile
	}
	fmt.Printf("magic = %v, version = %s, versionNum = %v\n", magic, version, versionNum)
	return versionNum, nil
}

func (app *application) sliceIndex(data []byte, sep byte) int {
	for i, b := range data {
		if b == sep {
			return i
		}
	}
	return -1
}
func (app *application) parseTable(bytes []byte) []byte {
	start := app.sliceIndex(bytes, resizeDBCode)
	end := app.sliceIndex(bytes, eofCode)
	return bytes[start+1 : end]
}
