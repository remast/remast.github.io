package main

import (
	"fmt"
	"time"
)

// START OMIT
func Hello(names chan string) {
	name := <-names // Nachricht empfangen // HL
	fmt.Println("Hello " + name + "!")
}

func main() {
	names := make(chan string) // Channel erzeugen // HL
	go Hello(names)
	names <- "Jim" // Nachricht senden // HL
	time.Sleep(time.Second)
}

// END OMIT
