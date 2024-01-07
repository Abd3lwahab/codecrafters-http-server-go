package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	okResponse := "HTTP/1.1 200 OK\r\n\r\n"
	notFoundResponse := "HTTP/1.1 404 Not Found\r\n\r\n"

	buffer := make([]byte, 1024)
	conn.Read(buffer)

	request := string(buffer)
	fmt.Println(request)

	path := strings.Split(request, " ")[1]

	if path == "/" {
		conn.Write([]byte(okResponse))
	} else {
		conn.Write([]byte(notFoundResponse))
	}
}
