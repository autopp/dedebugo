package main

import (
	"fmt"
	"os"
)

type Fmt struct{}

func (Fmt) Print(a ...interface{}) (int, error) {
	return 0, nil
}

func (Fmt) Println(a ...interface{}) (int, error) {
	return 0, nil
}

func (Fmt) Printf(format string, a ...interface{}) (int, error) {
	return 0, nil
}

func main() {
	fmt.Print("hello ")
	fmt.Println("world")
	fmt.Printf("from %s\n", os.Args[0])

	fmt := Fmt{}
	fmt.Print("hello ")
	fmt.Println("world")
	fmt.Printf("from %s\n", os.Args[0])
}
