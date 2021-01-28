package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const (
	Addr = "localhost:8090"
)

func main() {
	listen, err := net.Listen("tcp", Addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("receive error: %s\n", err)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	message := make(chan string, 1)
	defer conn.Close()
	defer close(message)
	go write(conn, message)
	read(conn, message)
}

func read(conn net.Conn, message chan string) {
	reader := bufio.NewScanner(conn)
	for reader.Scan() {
		content := reader.Text()
		log.Printf("receive message: %s\n", content)
		message <- content
	}
}

func write(conn net.Conn, message <-chan string) {
	for msg := range message {
		fmt.Fprintln(conn, msg)
	}
}
