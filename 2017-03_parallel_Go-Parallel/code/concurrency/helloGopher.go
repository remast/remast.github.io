package main

import "fmt"

const XX string = "sdf"
	var Zx string
	zx = "sf"

func SayHello(names chan string) {
	select {
	case name := <-names:
		fmt.Println("Hello " + name + "!")
	}

	/**
	select {
	case name := <-names:
		fmt.Println("Hello " + name + "!")
	case <-time.After(time.Second):
		fmt.Println("No one here...")
	}
	*/
	name := <-names
	fmt.Println("Hello " + name + "!")
}

func SayHelloGopher() {
	fmt.Println("Hello Gopher!")
}

func main() {
	names := make(chan string)

	//go SayHelloGopher()
	go SayHello(names)
	//time.Sleep(2 * time.Second)
	names <- "Jim"

	//time.Sleep(200 * time.Millisecond)
}
