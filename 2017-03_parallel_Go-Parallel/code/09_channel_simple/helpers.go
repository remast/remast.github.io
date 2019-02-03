// +build OMIT

package main

func main() {
	var value int

	// START1 OMIT
	// Deklaration und Initialisierung
	var c chan int
	c = make(chan int)
	// oder
	c := make(chan int) // HL
	// STOP1 OMIT

	// START2 OMIT
	// Nachricht senden
	c <- 1 // HL
	// STOP2 OMIT

	// START3 OMIT
	// Nachricht empfangen
	value = <-c // HL
	// STOP3 OMIT

	_ = value
}
