package main

import (
	"fmt"
	"time"
)

// START OMIT
func Hello(names chan string) {
	select { // HL
	case name := <-names: // HL
		fmt.Println("Hello " + name + "!")
	} // HL
}

func main() {
	names := make(chan string)
	go Hello(names)
	names <- "Jim"
	time.Sleep(time.Second)
}

// END OMIT
