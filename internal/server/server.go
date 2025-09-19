package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/response"
)

type Server struct {
	Port     int
	Listener net.Listener
	IsClosed atomic.Bool
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	server := &Server{
		Port:     port,
		Listener: l,
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
	err := response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		fmt.Println(err)
		return
	}
	h := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, h)
	if err != nil {
		fmt.Println(err)
		return
	}
}
