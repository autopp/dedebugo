package main

import (
	"fmt"
	"testing"
)

func TestSample(t *testing.T) {
	fmt.Print("hello ")
	fmt.Println("world")
	fmt.Printf("from %s\n", t.Name())

	fmt := Fmt{}
	fmt.Print("hello ")
	fmt.Println("world")
	fmt.Printf("from %s\n", t.Name())
}
