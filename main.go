package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
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
	f, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatalf("an error occured when openning the file : %s\n", err.Error())
	}

	ch := getLinesChannel(f)

	for str := range ch {
		fmt.Printf("read: %s\n", str)
	}
}
