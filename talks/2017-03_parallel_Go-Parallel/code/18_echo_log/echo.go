package main

import (
	"fmt"
	"io"
	"net"
)

// START 1 OMIT
func handler(c net.Conn, ch chan string) {
	ch <- c.RemoteAddr().String()
	io.Copy(c, c)
	defer c.Close()
}

func logger(ch chan string) {
	for {
		fmt.Println(<-ch)
	}
}

func server(l net.Listener, ch chan string) {
	for {
		c, err := l.Accept()
		if err != nil {
			continue
		}
		go handler(c, ch)
	}
}

// END 1 OMIT

// START 2 OMIT
func main() {
	l, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}
	messages := make(chan string)
	go logger(messages) // HL
	server(l, messages)
}

// END 2 OMIT
