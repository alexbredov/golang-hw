package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(line string) (string, error) {
	if len(line) == 0 {
		return "", nil
	}
	var result string
	runes := []rune(line)
	mem := runes[0]
	_, err := strconv.Atoi(string(mem))
	if err == nil {
		return result, ErrInvalidString
	}
	for i := 1; i < len(runes); i++ {
		current := runes[i]
		if _, err = strconv.Atoi(string(current)); err == nil {
			if _, err = strconv.Atoi(string(mem)); err == nil {
				return "", ErrInvalidString
			}
			number, _ := strconv.Atoi(string(current))
			result += strings.Repeat(string(mem), number)
			mem = current
			continue
		}
		if _, err = strconv.Atoi(string(mem)); err != nil {
			result += string(mem)
		}
		mem = current
	}
	if _, err = strconv.Atoi(string(mem)); err != nil {
		result += string(mem)
	}
	return result, nil
}
