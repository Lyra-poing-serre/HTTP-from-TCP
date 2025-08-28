package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/request"
)

func main() {
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Printf(
			"an error occured when setting up the TCP listener : %s\n",
			err.Error(),
		)
		return
	}
	defer ln.Close()
	fmt.Println("Listening for TCP traffic on :42069")

	for {
		con, err := ln.Accept()
		if err != nil {
			log.Printf(
				"an error occured when setting up the TCP listener : %s\n",
				err.Error(),
			)
			return
		}
		fmt.Println("A new connection has been accepted.")

		r, err := request.RequestFromReader(con)
		if err != nil {
			log.Printf(
				"an error occured while parsing the reqeust: %s\n",
				err.Error(),
			)
			return
		}

		fmt.Printf("Request line:\n- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)

		fmt.Println("The connection has been closed.")
	}
}
