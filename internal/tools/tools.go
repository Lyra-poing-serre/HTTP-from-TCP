package tools

import (
	"io"
	"unicode"
)

const CRLF = "\r\n"

type ChunkReader struct {
	Data            string
	NumBytesPerRead int
	Pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *ChunkReader) Read(p []byte) (n int, err error) {
	if cr.Pos >= len(cr.Data) {
		return 0, io.EOF
	}
	endIndex := cr.Pos + cr.NumBytesPerRead
	if endIndex > len(cr.Data) {
		endIndex = len(cr.Data)
	}
	n = copy(p, cr.Data[cr.Pos:endIndex])
	cr.Pos += n

	return n, nil
}

var allowedSpecialChars = map[rune]struct{}{ // ty Boots
	'!':  {},
	'#':  {},
	'$':  {},
	'%':  {},
	'&':  {},
	'\'': {}, // Note the escaping for the single quote
	'*':  {},
	'+':  {},
	'-':  {},
	'.':  {},
	'^':  {},
	'_':  {},
	'`':  {},
	'|':  {},
	'~':  {},
}

func IsForbiddenChar(r rune) bool {
	_, exist := allowedSpecialChars[r]
	return !unicode.IsNumber(r) && !unicode.IsLetter(r) && !exist
}
