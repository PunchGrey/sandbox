package main

import (
	"fmt"
)

func main() {
	t := "Hello"
	go func(s string) {
		fmt.Println("goroutine", s)
	}(t)
	fmt.Scanln()
}
