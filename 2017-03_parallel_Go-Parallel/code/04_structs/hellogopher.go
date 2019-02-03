package main

import "fmt"

// START OMIT
type Person struct { // HL
	Name string
} // HL

func (this Person) greet() {
	fmt.Println("Hello " + this.Name + "!")
}

func main() {
	p := Person{Name: "Gopher"} // HL
	p.greet()                   // HL
}

// END OMIT
