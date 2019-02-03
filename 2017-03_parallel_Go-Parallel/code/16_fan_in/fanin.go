package main

import (
	"fmt"
	"time"
)

// START OMIT
func producer(ch chan string, name string, d time.Duration) {
	var i int
	for {
		ch <- fmt.Sprintf("%s: %d", name, i)
		i++
		time.Sleep(d)
	}
}

func main() {
	ch := make(chan string)

	go producer(ch, "\t\tProducer 1", 200*time.Millisecond) // HL
	go producer(ch, "Producer 2", 800*time.Millisecond)     // HL

	for i := range ch { // Fan In // HL
		fmt.Println(i)
	} // HL
}

// END OMIT
