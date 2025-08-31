package headers

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Lyra-poing-serre/HTTP-from-TCP/internal/tools"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIdx := bytes.Index(data, []byte(tools.CRLF))
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
		if tools.IsForbiddenChar(r) {
			return 0, false, fmt.Errorf("invalid field-name: %s", key)
		}
	}
	err = h.Set(key, strData[keyIdx+1:])
	if err != nil {
		return 0, false, err
	}

	return crlfIdx + 2, false, nil
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
	return nil
}

func (h Headers) Get(key string) (string, error) {
	v, ok := h[strings.ToLower(key)]
	if !ok {
		return "", fmt.Errorf("key %s not found in headers", key)
	}
	return v, nil
}
