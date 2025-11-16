package protocol

import (
	"bufio"
	"strings"
)

// Decoder reads raw input lines from the client
// and turns them into parsed Command objects.
type Decoder struct {
	r *bufio.Reader
}

// Constructor
func NewDecoder(r *bufio.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode reads exactly one command line from the client.
func (d *Decoder) Decode() (Command, error) {
	line, err := d.r.ReadString('\n')
	if err != nil {
		// Any IO error is a system error, not a protocol error.
		return nil, err
	}

	// Trim CRLF / spaces
	line = strings.TrimSpace(line)

	if line == "" {
		return nil, ErrEmptyCommand
	}

	// Parse the line using your existing pipeline
	cmd, err := ParseLine(line)
	if err != nil {
		// If parser returned a protocol error, pass it through
		if IsProtocolError(err) {
			return nil, err
		}

		// If parser returned a generic error, wrap it
		return nil, NewProtocolError(err.Error())
	}

	return cmd, nil
}
