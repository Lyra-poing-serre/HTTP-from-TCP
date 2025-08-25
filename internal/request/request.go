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
		n, reqLine, err := parseRequestLine(string(data))
		if err != nil {
			return n, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = reqLine
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
	request := Request{
		State: Initialized,
	}
	buf := make([]byte, bufferSize)
	idx := 0

	for request.State != Done {
		if idx >= len(buf) {
			b := make([]byte, len(buf)*2)
			copy(b, buf)
			buf = b
		}
		n, err := reader.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.State = Done
				break
			}
			return &Request{}, err
		}
		idx += n

		n, err = request.parse(buf)
		if err != nil {
			return &Request{}, err
		}
		copy(make([]byte, bufferSize), buf)
		idx -= n
	}
	_, err := request.parse(buf)
	if err != nil {
		return &Request{}, err
	}
	return &request, nil
}

func parseRequestLine(request string) (int, RequestLine, error) {
	req := strings.Split(request, CRLF)
	// bytes.Index(data, []byte(CRLF))  // get first index
	if len(req) == 1 {
		return 0, RequestLine{}, nil
	}
	n := len(req)

	reqLine := strings.Split(req[0], " ")
	if len(reqLine) < 3 {
		return 0, RequestLine{}, fmt.Errorf("bad request-line format: %s", request)
	}
	method := reqLine[0]
	if method != GET && method != POST && method != PUT && method != DELETE {
		return 0, RequestLine{}, fmt.Errorf("invalid method: %s", method)
	}
	target := reqLine[1]
	http, ver, ok := strings.Cut(reqLine[2], "/")
	if !ok {
		return 0, RequestLine{}, fmt.Errorf("HTTP version not found: %s", reqLine[2])
	} else if http != "HTTP" {
		return 0, RequestLine{}, fmt.Errorf("unrecognized HTTP version: %s", reqLine[2])
	} else if ver != "1.1" {
		return 0, RequestLine{}, fmt.Errorf("only support HTTP/1.1: %s", reqLine[2])
	}

	return n, RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   ver,
	}, nil
}
