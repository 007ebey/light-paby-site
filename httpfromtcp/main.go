package main

import ("fmt"
         "os"
		 "io"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
      fmt.Println("Error:", err)
	  return
	}
	
	defer file.Close()

	buffer := make([]byte, 8)

	// Persistent variable
	var line string
	// infinite for loop
	for {
	   n, err := file.Read(buffer)
	   if err == io.EOF { 
		  fmt.Printf("read: %s\n", "end")
		  break
	   }
	   // Decoupled here, detecting nl
       part := -1
	   for i, b := range buffer[:n] {
		if b == '\n' {
            part = i + 1
		}
	   }

	   // there is a new line ( first line works great!)
	   if part > -1 {
		 // only update part of line
		 line += string(buffer[:part - 1])
		
		 fmt.Printf("read: %s\n", line)
		 // need to read the other part lol, forgot string!
		 line = string(buffer[part:n])
	   } else {
		 // keep updating the line
		 line += string(buffer[:n])
	   }

	   if err != nil {
		// oops forgot double quotes
		fmt.Println("Error:", err)
		break
	   }

	}
}