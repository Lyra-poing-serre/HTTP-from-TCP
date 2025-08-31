package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/headers"
	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/tools"
)

type ParseState int

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
)

const (
	parseInitialized ParseState = iota
	parseHeaders
	parseBody
	parseDone
)

const bufferSize = 8

type Request struct {
	RequestLine RequestLine
	State       ParseState
	Headers     headers.Headers
	Body        []byte
}

func NewRequest() *Request {
	return &Request{
		State:   parseInitialized,
		Headers: headers.NewHeaders(),
	}
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := NewRequest()
	buf := make([]byte, bufferSize)
	idx := 0

	for request.State != parseDone {
		if idx >= len(buf) {
			b := make([]byte, len(buf)*2)
			copy(b, buf)
			buf = b
		}
		n, err := reader.Read(buf[idx:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.State = parseDone
				break
			}
			return nil, err
		}
		idx += n
		n, err = request.parse(buf[:idx])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[n:])
		idx -= n
	}
	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case parseInitialized:
		n, requestLine, err := parseRequestLine(data)
		if err != nil {
			return n, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.State = parseHeaders
		return n, nil
	case parseHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return n, err
		}
		if n == 0 {
			return 0, nil
		}
		if done {
			r.State = parseBody
		}
		return n, nil
	case parseBody:

		return 0, nil
	case parseDone:
		return 0, errors.New("trying to read data in a done state")
	default:
		return 0, errors.New("unknown state")
	}
}

func parseRequestLine(data []byte) (int, *RequestLine, error) {
	idx := bytes.Index(data, []byte(tools.CRLF))
	if idx == -1 {
		return 0, nil, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return 0, nil, err
	}
	return idx + 2, requestLine, nil
}

func requestLineFromString(request string) (*RequestLine, error) {
	parts := strings.Split(request, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf(
			"bad request-line format: %s",
			request,
		)
	}
	method := parts[0]
	if method != GET && method != POST && method != PUT && method != DELETE {
		return nil, fmt.Errorf("invalid method: %s", method)
	}
	target := parts[1]
	http, ver, ok := strings.Cut(parts[2], "/")
	if !ok {
		return nil, fmt.Errorf(
			"HTTP version not found: %s",
			parts[2],
		)
	} else if http != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP version: %s", parts[2])
	} else if ver != "1.1" {
		return nil, fmt.Errorf("only support HTTP/1.1: %s", parts[2])
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   ver,
	}, nil
}
