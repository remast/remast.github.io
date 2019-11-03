package main

import (
	"fmt"
	"time"
)

// START OMIT

func player(name string, table chan int) {
	for {
		ball := <-table // Ball kommt // HL
		ball++
		fmt.Println(name, "schlägt", ball, "ten Ball")
		time.Sleep(400 * time.Millisecond)
		table <- ball // Ball zurückspielen // HL
	}
}

func main() {
	ball := 0
	table := make(chan int)

	go player("Jack", table) // Spieler 1 // HL
	go player("Tom", table)  // Spieler 2 // HL

	table <- ball
	time.Sleep(2 * time.Second) // Spielzeit // HL
	<-table
}

// END OMIT
