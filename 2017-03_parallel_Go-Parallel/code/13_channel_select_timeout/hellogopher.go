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
	case <-time.After(time.Second): // HL
		fmt.Println("Keiner da ...!")
	}
}

func main() {
	names := make(chan string)
	go Hello(names)
	time.Sleep(3 * time.Second)
	fmt.Println("Ich hau' ab.")
}

// END OMIT
