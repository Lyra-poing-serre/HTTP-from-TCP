package main

import (
	"log"
	"net"
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

}
