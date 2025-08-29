package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
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

const CRLF = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIdx := bytes.Index(data, []byte(CRLF))
	if crlfIdx == -1 {
		return 0, false, nil
	}
	if crlfIdx == 0 {
		return 2, true, nil
	}

	strData := string(data)[:crlfIdx]
	keyIdx := strings.Index(strData, ":")
	if keyIdx == -1 {
		return 0, false, fmt.Errorf("headers key not found: %s", strData)
	} else if strData[keyIdx-1] == ' ' {
		return 0, false, fmt.Errorf("malformed headers, no OWS next to ':' permitted: %s", strData[:keyIdx])
	}

	key := strings.ToLower(strings.TrimSpace(strData[:keyIdx]))
	for _, r := range key {
		if !unicode.IsNumber(r) &&
			!unicode.IsLetter(r) &&
			!isAllowedSpecialChar(r) {
			return 0, false, fmt.Errorf("invalid field-name: %s", key)
		}
	}
	err = h.Set(key, strData[keyIdx+1:])
	if err != nil {
		return 0, false, err
	}

	return crlfIdx + len(CRLF), false, nil
}

func (h Headers) Set(key, value string) error {
	valIdx := strings.LastIndex(value, ":")
	if valIdx != -1 {
		if value[valIdx-1] == ' ' || value[valIdx+1] == ' ' {
			return fmt.Errorf(
				"malformed headers, no OWS next to ':' permitted: %s",
				value,
			)
		}
		value = fmt.Sprintf(
			"%s%s",
			strings.TrimSpace(value[:valIdx]),
			strings.TrimSpace(value[valIdx:]),
		)
	} else {
		value = strings.TrimSpace(value)
	}
	_, exist := h[key]
	if !exist {
		h[key] = value
	} else {
		h[key] += fmt.Sprintf(", %s", value)
	}
	fmt.Printf("current key: %s -> %s\n", key, value)
	return nil
}

func isAllowedSpecialChar(r rune) bool {
	_, exist := allowedSpecialChars[r]
	return exist
}
