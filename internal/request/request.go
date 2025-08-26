package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type ParseState int

const CRLF = "\r\n"

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
)

const (
	Initialized ParseState = iota
	Done
)

const bufferSize = 8

type Request struct {
	RequestLine RequestLine
	State       ParseState
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case Initialized:
		reqLine, _, ok := strings.Cut(string(data), CRLF)
		if !ok {
			return 0, nil
		}

		n, requestLine, err := parseRequestLine(reqLine)
		if err != nil {
			return n, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = requestLine
		r.State = Done
		return n, nil
	case Done:
		return 0, errors.New("trying to read data in a done state")
	default:
		return 0, errors.New("unknown state")
	}
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{
		State: Initialized,
	}
	buf := make([]byte, bufferSize)
	idx := 0

	for request.State != Done {
		n, err := reader.Read(buf[idx:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.State = Done
				break
			}
			return &Request{}, err
		}
		idx += n
		if idx >= len(buf) {
			b := make([]byte, len(buf)*2)
			copy(b, buf)
			buf = b
		}

		n, err = request.parse(buf[:idx])
		if err != nil {
			return &Request{}, err
		}
		if n > 0 {
			copy(buf[:idx], make([]byte, bufferSize))
			idx -= n
		}
	}
	return request, nil
}

func parseRequestLine(request string) (int, RequestLine, error) {
	n := len(request)
	parts := strings.Split(request, " ")

	if len(parts) != 3 {
		return 0, RequestLine{}, fmt.Errorf(
			"bad request-line format: %s",
			request,
		)
	}
	method := parts[0]
	if method != GET && method != POST && method != PUT && method != DELETE {
		return 0, RequestLine{}, fmt.Errorf("invalid method: %s", method)
	}
	target := parts[1]
	http, ver, ok := strings.Cut(parts[2], "/")
	if !ok {
		return 0, RequestLine{}, fmt.Errorf(
			"HTTP version not found: %s",
			parts[2],
		)
	} else if http != "HTTP" {
		return 0, RequestLine{}, fmt.Errorf("unrecognized HTTP version: %s", parts[2])
	} else if ver != "1.1" {
		return 0, RequestLine{}, fmt.Errorf("only support HTTP/1.1: %s", parts[2])
	}

	return n, RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   ver,
	}, nil
}
