package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/request"
	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/response"
	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/tools"
)

type Server struct {
	Port        int
	Listener    net.Listener
	IsClosed    atomic.Bool
	HandlerFunc Handler
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode tools.StatusCode
	Message    string
}

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
	buf := bytes.NewBuffer([]byte{})

	req, err := request.RequestFromReader(conn)
	if err != nil {
		e := HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    err.Error(),
		}
		e.writeError(conn)
		return
	}
	req.Print()
	defer fmt.Println(
		"Server processed the request.\nWaiting another connection...",
	)

	e := s.HandlerFunc(buf, req)
	if e != nil {
		e.writeError(conn)
		return
	}

	body := buf.Bytes()
	err = response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		e := HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    err.Error(),
		}
		e.writeError(conn)
		return
	}
	h := response.GetDefaultHeaders(len(body))
	err = response.WriteHeaders(conn, h)
	if err != nil {
		e := HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    err.Error(),
		}
		e.writeError(conn)
		return
	}
	conn.Write(body)
}

func (e *HandlerError) writeError(w io.Writer) {
	err := response.WriteStatusLine(w, e.StatusCode)
	if err != nil {
		fmt.Println(err)
		return
	}
	h := response.GetDefaultHeaders(len(e.Message))
	err = response.WriteHeaders(w, h)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = w.Write([]byte(e.Message))
	if err != nil {
		fmt.Println(err)
		return
	}
}
