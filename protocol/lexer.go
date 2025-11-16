package protocol

import (
	"strings"
	"unicode"
)

// Token represents one piece of the command.
// Example: SET, foo, bar
type Token struct {
	Value string
}

// Lex takes a full command line and returns tokens.
// Only ASCII characters are allowed; whitespace and control characters are excluded.
func Lex(command string) ([]Token, error) {
	// Trim whitespace + newline
	line := strings.TrimSpace(command)
	if line == "" {
		return nil, ErrEmptyCommand
	}

	slic := strings.Fields(line) // Using Fields to handle multiple spaces and tabs
	tokens := make([]Token, 0, len(slic))

	for _, element := range slic {
		if element == "" {
			continue
		}

		// Validate ASCII + no control chars
		for _, r := range element {
			if r > 127 || unicode.IsControl(r) {
				return nil, NewProtocolError("invalid character in input")
			}
		}

		tokens = append(tokens, Token{Value: element})
	}

	return tokens, nil
}
