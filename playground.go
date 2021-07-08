package main

import "fmt"

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic: %+v\n", r)
		}
	}()
	A()
}

func A() {
	defer fmt.Println("D")
	panic("Error")
}
