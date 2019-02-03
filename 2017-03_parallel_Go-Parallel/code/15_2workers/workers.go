package main

import "fmt"
import "time"

// START OMIT
func worker(name string, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		time.Sleep(time.Second)
		fmt.Println(name, "ist fertig mit Job", j)
		results <- j
	}
}

func main() {
	jobs := make(chan int, 10)    // Channel für Jobs mit Puffer Größe 10 // HL
	results := make(chan int, 10) // Channel für Ergebnisse mit Puffer Größe 10 // HL

	go worker("John", jobs, results)  // Worker 1 // HL
	go worker("Maria", jobs, results) // Worker 2 // HL

	for j := 1; j <= 5; j++ { // 5 Jobs einstellen
		jobs <- j // HL
	}
	close(jobs) // Channel schließen // HL

	for a := 1; a <= 5; a++ { // Ergebnisse abholen
		fmt.Println("Ergebnis ist", <-results) // HL
	}
}

// END OMIT
