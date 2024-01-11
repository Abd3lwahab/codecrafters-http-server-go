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
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	notFoundResponse := "HTTP/1.1 404 Not Found\r\n\r\n"

	buffer := make([]byte, 1024)
	conn.Read(buffer)

	request := string(buffer)
	fmt.Println(request)

	path := strings.Split(request, " ")[1]
	fmt.Println(path)

	if path == "/" {
		conn.Write([]byte(okResponse("")))
	} else if strings.HasPrefix(path, "/echo/") {
		body := strings.TrimPrefix(path, "/echo/")
		conn.Write([]byte(okResponse(body)))
	} else if path == "/user-agent" {
		userAgent := strings.Split(strings.Split(request, "\r\n")[2], " ")[1]
		conn.Write([]byte(okResponse(userAgent)))
	} else {
		conn.Write([]byte(notFoundResponse))
	}
}

func okResponse(body string) string {
	if body == "" {
		return "HTTP/1.1 200 OK\r\n\r\n"
	}
	headers := "Content-Type: text/plain\r\n"

	return "HTTP/1.1 200 OK\r\n" + headers + "Content-Length: " + fmt.Sprint(len(body)) + "\r\n\r\n" + body
}
