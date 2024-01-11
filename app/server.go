package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

var dirFlag = flag.String("directory", ".", "Directory to serve")

func main() {
	fmt.Println("Logs from your program will appear here!")

	flag.Parse()
	fmt.Println("Serving directory: " + *dirFlag)

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
		conn.Write([]byte(okResponse("", "")))
	} else if strings.HasPrefix(path, "/echo/") {
		body := strings.TrimPrefix(path, "/echo/")
		conn.Write([]byte(okResponse(body, "text/plain")))
	} else if path == "/user-agent" {
		userAgent := strings.Split(strings.Split(request, "\r\n")[2], " ")[1]
		conn.Write([]byte(okResponse(userAgent, "text/plain")))
	} else if strings.HasPrefix(path, "/files/") {
		fileName := strings.TrimPrefix(path, "/files/")
		filePath := filepath.Join(*dirFlag, fileName)

		fmt.Println(filePath)

		if _, err := os.Stat(filePath); err != nil {
			conn.Write([]byte(notFoundResponse))
			return
		}

		_, err := os.Open(filePath)
		if err != nil {
			conn.Write([]byte(notFoundResponse))
			return
		}

		buffer, err := os.ReadFile(filePath)
		if err != nil {
			conn.Write([]byte(notFoundResponse))
			return
		}

		_, err = os.Stdout.Write(buffer)
		if err != nil {
			conn.Write([]byte(notFoundResponse))
			return
		}

		fileData := string(buffer)
		fmt.Println(fileData)

		conn.Write([]byte(okResponse(fileData, "application/octet-stream")))
	} else {
		conn.Write([]byte(notFoundResponse))
	}
}

func okResponse(body string, contentType string) string {
	if body == "" {
		return "HTTP/1.1 200 OK\r\n\r\n"
	}
	headers := fmt.Sprintf("Content-Type: %s\r\n", contentType)

	return "HTTP/1.1 200 OK\r\n" + headers + "Content-Length: " + fmt.Sprint(len(body)) + "\r\n\r\n" + body
}
