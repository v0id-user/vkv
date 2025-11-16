package main

import "fmt"

func main () {
	go fmt.Println("Hello, World!")
	go fmt.Println("Hello, World!")
	go fmt.Println("Hello, World!")

	select {}
}