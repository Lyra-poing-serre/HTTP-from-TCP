package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		buf := make([]byte, 8)
		line := ""
		for {
			n, err := f.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				log.Fatalf("an error occured while reading the file : %s\n", err.Error())
			}
			s := strings.Split(string(buf[:n]), "\n")
			line += s[0]
			if len(s) > 1 {
				ch <- line
				line = strings.Join(s[1:], "")
			}
		}
		ch <- line
	}()

	return ch
}

func main() {
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("an error occured when setting up the TCP listener : %s\n", err.Error())
	}
	defer ln.Close()
	fmt.Println("Listening for TCP traffic on :42069")

	for {
		con, err := ln.Accept()
		if err != nil {
			log.Fatalf("an error occured when setting up the TCP listener : %s\n", err.Error())
		}
		fmt.Println("A new connection has been accepted.")

		ch := getLinesChannel(con)

		for str := range ch {
			fmt.Println(str)
		}
		fmt.Println("The connection has been closed.")
	}
}
