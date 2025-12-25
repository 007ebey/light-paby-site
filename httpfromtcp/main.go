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
	defer file.Close()
	for line := range getLinesChannel(file) {
		fmt.Println(line)
	}
}

func getLinesChannel(f io.ReadCloser) <- chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		defer f.Close()
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
