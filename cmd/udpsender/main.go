package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	address, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("an error occured when resolving the address: %s\n", err.Error())
	}
	con, err := net.DialUDP("udp", nil, address)
	if err != nil {
		log.Fatalf("an error occured when setting up the UDP connection: %s\n", err.Error())
	}
	defer con.Close()

	bufReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		str, err := bufReader.ReadString(byte('\n'))
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("End detected, leaving for loop.")
				break
			} else {
				log.Fatalf("an error occured when setting up the UDP connection: %s\n", err.Error())
			}
		}
		_, err = con.Write([]byte(str))
		if err != nil {
			log.Fatalf("an error occured when writing to the UDP connection: %s\n", err.Error())
		}
	}

}
