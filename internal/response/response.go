package response

import (
	"errors"
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
	//
	WriterStatusLine tools.WriterState = 0
	WriterHeaders    tools.WriterState = 1
	WriterBoby       tools.WriterState = 2
)

type Writer struct {
	writerState tools.WriterState
	Connection  io.Writer
}

func (w *Writer) Write(b []byte) (int, error) {
	return w.Connection.Write(b)
}

func (w *Writer) WriteStatusLine(statusCode tools.StatusCode) error {
	if w.writerState != WriterStatusLine {
		return errors.New("writer not in status line states")
	}
	var err error
	switch statusCode {
	case 200:
		// fmt.FprintF remplace w.Write([]byte(fmt.Sprintf(...))
		_, err = fmt.Fprintf(
			w, "HTTP/1.1 %d %s%s",
			statusCode, "OK", tools.CRLF,
		)
	case 400:
		_, err = fmt.Fprintf(
			w, "HTTP/1.1 %d %s%s",
			statusCode, "Bad Request", tools.CRLF,
		)
	case 500:
		_, err = fmt.Fprintf(
			w, "HTTP/1.1 %d %s%s",
			statusCode, "Internal Server Error", tools.CRLF,
		)
	default:
		_, err = fmt.Fprintf(w, "HTTP/1.1 %d %s", statusCode, tools.CRLF)
	}
	if err == nil {
		w.writerState = WriterHeaders
	}
	return err
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
	err = h.Set("Content-Type", "text/html")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return h
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != WriterHeaders {
		return errors.New("writer not in headers states")
	}
	b := []byte{}
	for k, v := range headers {
		b = fmt.Appendf(b, "%s: %s%s", k, v, tools.CRLF)
	}
	b = append(b, tools.CRLF...)
	_, err := w.Write(b)
	if err == nil {
		w.writerState = WriterBoby
	}
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != WriterBoby {
		return 0, errors.New("writer not in body states")
	}
	return w.Write(p)
}
