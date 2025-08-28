package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with : in host name
	example := "Ghost: https://google.com:69\r\n"
	headers = NewHeaders()
	data = []byte(example + "\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "https://google.com:69", headers["ghost"])
	assert.Equal(t, len(example), n)
	assert.False(t, done)

	// Test: Valid single header with space
	example = "Host:                        localhost:42069                 \r\n"
	headers = NewHeaders()
	data = []byte(example + "\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, len(example), n)
	assert.False(t, done)

	// Test: Valid double headers with a done check
	example = "Ghost: https://google.com:69\r\n"
	ex := "Host: localhost:42069\r\n"
	headers = NewHeaders()
	data = []byte(example + ex + "\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "https://google.com:69", headers["ghost"])
	assert.Equal(t, len(example), n)
	assert.False(t, done)
	data = data[n:]
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "https://google.com:69", headers["ghost"])
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.False(t, done)
	_, done, err = headers.Parse(data[n:])
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "https://google.com:69", headers["ghost"])
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid single header with space
	example = "Host: localhost :42069\r\n"
	headers = NewHeaders()
	data = []byte(example + "\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid field-name header
	headers = NewHeaders()
	data = []byte("H@st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
