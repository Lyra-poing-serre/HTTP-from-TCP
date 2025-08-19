package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatalf("an error occured when openning the file : %s\n", err.Error())
	}
	defer f.Close()

	buf := make([]byte, 8)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("an error occured while reading the file : %s\n", err.Error())
		}
		fmt.Printf("read: %s\n", string(buf[:n]))
	}
}
