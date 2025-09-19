package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/headers"
	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/tools"
)

const (
	// Pour référence, le enum peux aussi se faire comme cela, avec _ pour skip
	// _ StatusCode = iota * 100
	// _
	// StatusOK
	// _
	// StatusBadRequest
	// StatusInternalServerError
	StatusOK                  tools.StatusCode = 200
	StatusBadRequest          tools.StatusCode = 400
	StatusInternalServerError tools.StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode tools.StatusCode) error {
	switch statusCode {
	case 200:
		// fmt.FprintF remplace w.Write([]byte(fmt.Sprintf(...))
		_, err := fmt.Fprintf(
			w, "HTTP/1.1 %d %s%s",
			statusCode, "OK", tools.CRLF,
		)
		return err
	case 400:
		_, err := fmt.Fprintf(
			w, "HTTP/1.1 %d %s%s",
			statusCode, "Bad Request", tools.CRLF,
		)
		return err
	case 500:
		_, err := fmt.Fprintf(
			w, "HTTP/1.1 %d %s%s",
			statusCode, "Internal Server Error", tools.CRLF,
		)
		return err
	default:
		_, err := fmt.Fprintf(w, "HTTP/1.1 %d %s", statusCode, tools.CRLF)
		return err
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	err := h.Set("Content-Length", strconv.Itoa(contentLen))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	err = h.Set("Connection", "close")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	err = h.Set("Content-Type", "text/plain")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	b := []byte{}
	for k, v := range headers {
		b = fmt.Appendf(b, "%s: %s%s", k, v, tools.CRLF)
	}
	b = append(b, tools.CRLF...)
	_, err := w.Write(b)
	return err
}
