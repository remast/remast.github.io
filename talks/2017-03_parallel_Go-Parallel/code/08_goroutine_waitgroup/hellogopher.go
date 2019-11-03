package main

import (
	"fmt"
	"sync"
)

// START OMIT
func HelloGopher() {
	fmt.Println("Hello Gopher!")
}

func main() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go func() { // HL
		defer waitGroup.Done()
		HelloGopher()
	}() // HL
	waitGroup.Wait()
}

// END OMIT
