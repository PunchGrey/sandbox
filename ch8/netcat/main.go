package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	ports := os.Args[1:]
	if len(ports) == 0 {
		ports = append(ports, "8000")
	}
	fmt.Println("test", ports)
	/*conn, _ := net.Dial("tcp", "localhost:"+ports[0])
	mustCopy(os.Stdout, conn)
	defer conn.Close()*/
	chPort := make(chan string, 1)
	for _, port := range ports {
		fmt.Println("test", port)
		chPort <- port
		go getTime(chPort)
	}

	fmt.Scanln()
	/*
		ports := flag.Args()
		if ports == nil {
			ports = append(ports, "8001")
		}
		chanPort := make(chan string)
		//	chWrite := make(chan io.Writer)
		//chWrite <- os.Stdout

		for _, port := range ports {
			fmt.Println("localhost:" + port)

			func(ch <-chan string) {
				fmt.Println("befor ch")
				port := <-ch
				fmt.Println("after ch", port)
				conn, err := net.Dial("tcp", "localhost:"+port)
				fmt.Println("erriconn")
				if err != nil {
					log.Fatal(err)
				}
				defer conn.Close()

				mustCopy(os.Stdout, conn)
			}(chanPort)
			chanPort <- port
		}
		close(chanPort)*/
}

func mustCopy(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	if err != nil {
		log.Fatal(err)
	}
}

func getTime(ch <-chan string) {
	port := <-ch
	fmt.Println("test", "localhost:"+port)
	conn, err := net.Dial("tcp", "localhost:"+port)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}
	mustCopy(os.Stdout, conn)
	fmt.Println("test", conn)
}
