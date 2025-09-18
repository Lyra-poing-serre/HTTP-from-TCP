package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	Port     int
	Listener net.Listener
	IsClosed atomic.Bool
}

const response = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!"

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
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Printf("an error occured when writing response: %s", err.Error())
	}
}
