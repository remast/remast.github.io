package main

import (
	"fmt"
	"time"
)

// START OMIT
func Hello(names chan string) {
	select {
	case name := <-names:
		fmt.Println("Hello " + name + "!")
	default: // auf keinem Channels kommt eine Nachricht // HL
		fmt.Println("Keiner da...")
	}
}

func main() {
	names := make(chan string)
	go Hello(names)
	time.Sleep(time.Second)
}

// END OMIT
