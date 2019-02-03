package main

import "fmt"

// START OMIT

type Person struct {
	Name string
}

func nameFritz(p *Person) { // Pointer auf Person // HL
	p.Name = "Fritz"
}

func main() {
	p := Person{Name: "Max"}
	nameFritz(&p) // Referenz auf p übergeben // HL
	fmt.Println("Hello " + p.Name)
}

// END OMIT
