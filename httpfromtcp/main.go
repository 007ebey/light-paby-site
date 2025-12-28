package main

import (
  "fmt"
  "io"
  "net"
)

func main() {
	ln, _ := net.Listen("tcp", ":42069")
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		fmt.Println("Conection has been accepted")
		if err != nil {
          fmt.Println(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	for line := range getLinesChannel(conn) {
		fmt.Println(line)
	}
}

func getLinesChannel(f io.ReadCloser) <- chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		var line []byte
		buffer := make([]byte, 8)
		for {
			n, err := f.Read(buffer)
			for _, b := range buffer[:n] {
				if b == '\n' {
	                ch <- "read: " + string(line)
                    line = line[:0]
				} else {
					line = append(line, b)
				}
			}

			if err == io.EOF {
			    break
		    }

			if err != nil {
				return
			}
		}

		if len(line) > 0 {
			ch <- "read: " + string(line)
		}
	}()

	return ch
}
