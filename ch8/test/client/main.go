package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8050")
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}
	mustCopy(os.Stdout, conn)
	fmt.Println("test", conn)
}

func mustCopy(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	if err != nil {
		log.Fatal(err)
	}
}
