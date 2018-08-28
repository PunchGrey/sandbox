package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	fmt.Println("Hello")
	port := flag.String("p", "8000", "This is a port of server-time")
	timeZone := flag.String("tz", "Europe/Moscow", "this is a time zone")
	flag.Parse()

	//	tz := make(chan string, 1)
	//	tz <- *timeZone

	fmt.Println("Hello,", *port, *timeZone)
	listener, err := net.Listen("tcp", "localhost:"+*port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		fmt.Println("test", conn)
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn, *timeZone)
	}
}

func handleConn(c net.Conn, timeZone string) {
	defer c.Close()
	loc, _ := time.LoadLocation(timeZone)
	for {
		_, err := io.WriteString(c, timeZone+"\t"+time.Now().In(loc).Format("15:04:05\n"))
		if err != nil {
			defer c.Close()
			return
		}
		time.Sleep(1 * time.Second)
	}
}
