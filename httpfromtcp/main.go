package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
		}
		for line := range getLinesChannel(conn) {
			fmt.Println(line)
			break
		}
		conn.Close()
	}

}

// <-chan string is read only channel, that sends strings
func getLinesChannel(f net.Conn) <-chan string {
	ch := make(chan string)
	// reading loop inside a go routine thread (lightweight)
	go func() {
		defer close(ch)
		buffer := make([]byte, 8)
		var line string
		for {
			n, err := f.Read(buffer)

			part := -1
			for i, b := range buffer[:n] {
				if b == '\n' {
					part = i + 1
				}
			}

			if part > -1 {
				line += string(buffer[:part-1])
				ch <- line // send the new line to channel
				line = string(buffer[part:n])
			} else {
				line += string(buffer[:n])
			}

			if err == io.EOF {
				break
			}

			if err != nil {
				break
			}
		}
	}()
	return ch
}
