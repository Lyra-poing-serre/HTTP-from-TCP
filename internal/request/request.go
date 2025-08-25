package request

import (
	"fmt"
	"io"
	"strings"
)

const CRLF = "\r\n"
const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req, err := io.ReadAll(reader)
	if err != nil {
		return &Request{}, err
	}
	reqLine, err := parseRequestLine(string(req))
	if err != nil {
		return &Request{}, err
	}
	return &Request{
		RequestLine: reqLine,
	}, nil
}

func parseRequestLine(request string) (RequestLine, error) {
	req := strings.Split(request, CRLF)
	// bytes.Index(data, []byte(CRLF))  // get first index
	reqLine := strings.Split(req[0], " ")
	if len(reqLine) != 3 {
		return RequestLine{}, fmt.Errorf("bad request-line format: %s", request)
	}
	method := reqLine[0]
	if method != GET && method != POST && method != PUT && method != DELETE {
		return RequestLine{}, fmt.Errorf("invalid method: %s", method)
	}
	target := reqLine[1]
	http, ver, ok := strings.Cut(reqLine[2], "/")
	if !ok {
		return RequestLine{}, fmt.Errorf("HTTP version not found: %s", reqLine[2])
	} else if http != "HTTP" {
		return RequestLine{}, fmt.Errorf("unrecognized HTTP version: %s", reqLine[2])
	} else if ver != "1.1" {
		return RequestLine{}, fmt.Errorf("only support HTTP/1.1: %s", reqLine[2])
	}

	return RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   ver,
	}, nil
}
