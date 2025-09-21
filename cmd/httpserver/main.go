package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/request"
	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/response"
	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/server"
)

const (
	port           = 42069
	badRequestHTML = `
<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`
	internalServerError = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
	okResponse = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
)

func main() {
	server, err := server.Serve(
		port,
		func(w response.Writer, req *request.Request) {
			switch req.RequestLine.RequestTarget {
			case "/yourproblem":
				err := w.WriteStatusLine(response.StatusBadRequest)
				if err != nil {
					fmt.Println(err)
					return
				}
				h := response.GetDefaultHeaders(len(badRequestHTML))
				err = w.WriteHeaders(h)
				if err != nil {
					fmt.Println(err)
					return
				}
				_, err = w.WriteBody([]byte(badRequestHTML))
				if err != nil {
					fmt.Println(err)
					return
				}
				return
			case "/myproblem":
				err := w.WriteStatusLine(response.StatusInternalServerError)
				if err != nil {
					fmt.Println(err)
					return
				}
				h := response.GetDefaultHeaders(len(internalServerError))
				err = w.WriteHeaders(h)
				if err != nil {
					fmt.Println(err)
					return
				}
				_, err = w.WriteBody([]byte(internalServerError))
				if err != nil {
					fmt.Println(err)
					return
				}
				return
			default:
				err := w.WriteStatusLine(response.StatusOK)
				if err != nil {
					fmt.Println(err)
					return
				}
				h := response.GetDefaultHeaders(len(okResponse))
				err = w.WriteHeaders(h)
				if err != nil {
					fmt.Println(err)
					return
				}
				_, err = w.WriteBody([]byte(okResponse))
				if err != nil {
					fmt.Println(err)
					return
				}
				return
			}
		},
	)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
