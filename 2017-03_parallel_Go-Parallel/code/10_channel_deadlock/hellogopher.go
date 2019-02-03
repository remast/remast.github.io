package main

import "fmt"

// START OMIT
func Hello(names chan string) {
	name := <-names
	fmt.Println("Hello " + name + "!")
}

func main() {
	names := make(chan string)
	go Hello(names)
	names <- "Jim"
	names <- "Jack" // Nachricht senden // HL
}

// END OMIT
