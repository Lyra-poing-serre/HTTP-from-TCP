package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/request"
	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/response"
)

type Server struct {
	Port        int
	Listener    net.Listener
	IsClosed    atomic.Bool
	HandlerFunc Handler
}

type Handler func(w response.Writer, req *request.Request)

// type HandlerError struct {
// 	StatusCode tools.StatusCode
// 	Message    string
// }

func Serve(port int, h Handler) (*Server, error) {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	server := &Server{
		Port:        port,
		Listener:    l,
		HandlerFunc: h,
	}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.IsClosed.Store(true)
	err := s.Listener.Close()
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}

func (s *Server) listen() {
	for !s.IsClosed.Load() {
		conn, err := s.Listener.Accept()
		fmt.Println("New connection !")
		if err != nil && !s.IsClosed.Load() {
			fmt.Println(err)
			break
		}
		if conn != nil {
			go s.handle(conn)
		}
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	responseWriter := response.Writer{
		Connection: conn,
	}

	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Print()
	defer fmt.Println(
		"Server processed the request.\nWaiting another connection...",
	)

	s.HandlerFunc(responseWriter, req)
}
