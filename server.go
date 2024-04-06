package http

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func Listen() error {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Failed to read from connection: %v\n", err)
			return
		}
		if errors.Is(err, io.EOF) || n == 0 {
			continue
		}

		msg := ParseMessage(buf[:n])

		if msg.Path == "/" {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		} else if strings.HasPrefix(msg.Path, "/echo") {
			content := strings.TrimPrefix(msg.Path, "/echo/")
			conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(content), content)))
		} else if msg.Path == "/user-agent" {
			content := msg.Headers["User-Agent"]
			conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(content), content)))
		} else {
			conn.Write([]byte("HTTP/1.1 404 NOT FOUND\r\n\r\n"))
		}
	}
}

type HttpMessage struct {
	Method, Path, Version string
	Headers               map[string]string
}

func ParseMessage(raw []byte) HttpMessage {
	requestArray := strings.Split(string(raw), "\r\n")
	startLine := strings.Split(requestArray[0], " ")
	message := HttpMessage{
		Method:  startLine[0],
		Path:    startLine[1],
		Version: startLine[2],
		Headers: make(map[string]string),
	}
	for _, line := range requestArray[1:] {
		if line == "" {
			break
		}
		tmp := strings.Split(line, ": ")
		message.Headers[tmp[0]] = tmp[1]
	}
	return message
}
