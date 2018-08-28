package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8050")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		defer conn.Close()
		fmt.Println("test", conn)
		if err != nil {
			log.Print(err)
			continue
		}
		for {
			_, err := io.WriteString(conn, time.Now().Format("15:04:05\n"))
			if err != nil {
				conn.Close()
				break
			}
			time.Sleep(1 * time.Second)
		}
	}
}
