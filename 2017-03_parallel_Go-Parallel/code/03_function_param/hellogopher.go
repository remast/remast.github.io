package main

import "fmt"

// START OMIT
func Hello(name string) string { // HL
	return "Hello " + name + "!"
} // HL

func main() {
	fmt.Println(Hello("Gopher"))
}

// END OMIT
