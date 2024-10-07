package main
import (
	"fmt"
	"golang.org/x/example/hello/reverse"
)

func main() {
	initString := "Hello, OTUS!"
	fmt.Println(reverse.String(initString))
}
