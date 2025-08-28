package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

const CRLF = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	strData := string(data)
	fmt.Printf("got this data: %s\n", strData)
	crlfIdx := strings.Index(strData, CRLF)
	if crlfIdx == -1 {
		return 0, false, nil
	}
	if crlfIdx == 0 {
		return len(data), true, nil
	}

	strData = strData[:crlfIdx]
	keyIdx := strings.Index(strData, ":")
	if keyIdx == -1 {
		return 0, false, fmt.Errorf("headers key not found: %s", strData)
	} else if strData[keyIdx-1] == ' ' {
		return 0, false, fmt.Errorf("malformed headers, no OWS next to ':' permitted: %s", strData[:keyIdx])
	}
	key := strings.TrimSpace(strData[:keyIdx])
	value := strData[keyIdx+1:]

	valIdx := strings.LastIndex(value, ":")
	if valIdx == -1 {
		return 0, false, fmt.Errorf(
			"headers value not found: %s",
			strData[keyIdx:],
		)
	} else if value[valIdx-1] == ' ' || value[valIdx+1] == ' ' {
		return 0, false, fmt.Errorf("malformed headers, no OWS next to ':' permitted: %s", value)
	}
	h[key] = fmt.Sprintf(
		"%s%s",
		strings.TrimSpace(value[:valIdx]),
		strings.TrimSpace(value[valIdx:]),
	)
	return len(strData) + len(CRLF), false, nil
}

func NewHeaders() Headers {
	return Headers{}
}
