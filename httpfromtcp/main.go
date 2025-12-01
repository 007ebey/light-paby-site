package main

import ("fmt"
         "os"
		 "io"
)

func main() {
	fmt.Println("I hope I get the job!")
	file, err := os.Open("messages.txt")
	if err != nil {
      fmt.Println("Error:", err)
	  return
	}
	// close as per pleasure
	defer file.Close()
	// we need a buffer, slice of 8 bytes
	buffer := make([]byte, 8)

	// infinite for loop
	for {
       /*
	     number of bytes actually read, error
		 if any, data stored as reference
	   */ 
	   // Step 1 completed
	   n, err := file.Read(buffer)
	   // Step 3 completed
	   if err == io.EOF {
		  // no semicolon required
		  break
	   }
	   // Step 2, slice the characters read
	   toPrint := string(buffer[:n])
	   fmt.Printf("read: %s\n", toPrint)


	   if err != nil {
		// oops forgot double quotes
		fmt.Println("Error:", err)
		break
	   }

	}
}