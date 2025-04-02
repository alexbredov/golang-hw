package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrReadDir = errors.New("error reading directory")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	var fileName string
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	dirEntries, err := os.ReadDir(dir)
	envDirMap := make(Environment)
	if err != nil {
		return nil, ErrReadDir
	}
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			fileName = entry.Name()
			envName := strings.ReplaceAll(fileName, "=", "")
			fileInfo, err := os.Stat(dir + fileName)
			if err != nil {
				fmt.Printf("failed to read info of %v file: %v\n", fileName, err)
				continue
			}
			if fileInfo.Size() == 0 {
				envDirMap[envName] = EnvValue{Value: "", NeedRemove: true}
				continue
			}
			file, err := os.Open(dir + fileName)
			if err != nil {
				fmt.Printf("failed to open file %v: %v\n", fileName, err)
				continue
			}
			line, err := readLine(file)
			if err != nil {
				fmt.Printf("failed to read file %v: %v\n", fileName, err)
				continue
			}
			line = strings.TrimRight(line, "\t")
			envDirMap[envName] = EnvValue{Value: line, NeedRemove: false}
			file.Close()
		}
	}
	return envDirMap, nil
}

func readLine(reader io.Reader) (line string, err error) {
	var n int
	var lineSlice []byte
	runeSlice := make([]byte, 1)
	for {
		n, err = reader.Read(runeSlice)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if n > 0 {
			char := runeSlice[0]
			if char == '\n' {
				break
			}
			if char == 0x00 {
				char = '\n'
			}
			lineSlice = append(lineSlice, char)
		}
	}
	return string(lineSlice), nil
}
