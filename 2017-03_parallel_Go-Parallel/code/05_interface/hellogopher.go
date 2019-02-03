package main

import "fmt"

type Stepmother struct{}

func (p Stepmother) greet() {
	fmt.Println("Go to hell!")
}

// START OMIT
type Person struct {
	Name string
}

func (p Person) greet() {
	fmt.Println("Hello " + p.Name + "!")
}

type NicePerson interface { // HL
	greet() // HL
} // HL

func passBy(p1 NicePerson, p2 NicePerson) { // HL
	p1.greet()
	p2.greet()
} // HL

func main() {
	p := Person{Name: "Gopher"}
	s := Person{Name: "Parallel"}
	passBy(p, s) // HL
}

// END OMIT
