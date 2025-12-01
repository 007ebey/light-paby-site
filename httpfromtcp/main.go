package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer file.Close()

	// seperation of concern, reading part of complex file
	for line := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", line)
	}

}

// <-chan string is read only channel, that sends strings
func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	// reading loop inside a go routine thread (lightweight)
	go func() {
		defer close(ch)
		buffer := make([]byte, 8)
		var line string
		for {
			n, err := f.Read(buffer)
			if err == io.EOF {
				ch <- "end"
				return
			}
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

			if err != nil {
				// oops forgot double quotes
				fmt.Println("Error:", err)
				break
			}
		}
	}()
	return ch
}
