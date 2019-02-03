package main

import "fmt"

// START OMIT
func HelloGopher() {
	fmt.Println("Hello Gopher!")
}

func main() {
	go HelloGopher() // HL
}

// END OMIT
