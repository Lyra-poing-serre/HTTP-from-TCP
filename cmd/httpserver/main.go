package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/headers"
	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/request"
	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/response"
	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/server"
	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/tools"
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
			if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
				// ctx timeout todo
				chHttpbin := make(chan []byte)
				target := fmt.Sprintf(
					"%s%s",
					"https://httpbin.org",
					strings.TrimPrefix(
						req.RequestLine.RequestTarget,
						"/httpbin",
					),
				)
				go func() {
					resp, err := http.Get(target)
					if err != nil {
						fmt.Println(err)
						return
					}
					defer close(chHttpbin)
					defer resp.Body.Close()
					for {
						b := make([]byte, tools.ChunkSize)
						n, err := resp.Body.Read(b)
						if n > 0 {
							chHttpbin <- b[:n]
						}
						if errors.Is(err, io.EOF) {
							return
						}
					}
				}()

				err := w.WriteStatusLine(response.StatusOK)
				if err != nil {
					fmt.Println(err)
					return
				}
				h := headers.NewHeaders()
				h.Set("Content-Type", "application/json")
				h.Set("Transfer-Encoding", "chunked")
				h.Set("Trailer", "X-Content-SHA256")
				h.Set("Trailer", "X-Content-Length")
				err = w.WriteHeaders(h)
				if err != nil {
					fmt.Println(err)
					return
				}

				body := []byte{}
				for chunk := range chHttpbin {
					body = append(body, chunk...)
					_, err = w.WriteChunkedBody(chunk)
					if err != nil {
						fmt.Println(err)
						return
					}
				}
				_, err = w.WriteChunkedBodyDone()
				if err != nil {
					fmt.Println(err)
					return
				}

				sum := sha256.Sum256(body)
				h.Set("X-Content-SHA256", hex.EncodeToString(sum[:]))
				h.Set("X-Content-Length", strconv.Itoa(len(body)))
				err = w.WriteTrailers(h)
				if err != nil {
					fmt.Println(err)
				}
				return
			} else {
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
