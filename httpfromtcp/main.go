package main

import (
  "fmt"
  "os"
  "io"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Failed to open")
	}
	buff := make([]byte, 8)
	defer file.Close()
	for {
		_, err := file.Read(buff)
		fmt.Printf("read: %s\n", buff)

		if err == io.EOF {
			break
		}
	}
}
