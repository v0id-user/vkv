package protocol

import (
	"strings"
)

func ParseTokens(tokens []Token) (Command, error) {
	if len(tokens) == 0 {
		return nil, NewProtocolError("no tokens to parse")
	}

	// Operation
	op := strings.ToUpper(tokens[0].Value)

	switch op {
	case "SET":
		// SET key value
		if len(tokens) != 3 {
			return nil, Errorf("SET requires exactly 2 arguments: key and value, got %d", len(tokens)-1)
		}
		key := tokens[1].Value
		val := tokens[2].Value

		return Set{
			Key:   key,
			Value: val,
		}, nil
	case "GET":
		// GET key
		if len(tokens) != 2 {
			return nil, Errorf("GET requires exactly 1 argument: key, got %d", len(tokens)-1)
		}
		key := tokens[1].Value
		if key == "" {
			return nil, NewProtocolError("GET key cannot be empty")
		}

		return Get{
			Key: key,
		}, nil
	case "DEL":
		// DEL key
		if len(tokens) != 2 {
			return nil, Errorf("DEL requires exactly 1 argument: key, got %d", len(tokens)-1)
		}
		key := tokens[1].Value
		if key == "" {
			return nil, NewProtocolError("DEL key cannot be empty")
		}

		return Del{
			Key: key,
		}, nil
	case "FLUSH":
		// FLUSH command takes no arguments
		if len(tokens) != 1 {
			return nil, Errorf("FLUSH does not take any arguments, got %d", len(tokens)-1)
		}
		return Flush{}, nil
	default:
		return nil, Errorf("unknown command: %q", tokens[0].Value)
	}
}

// ParseLine is a helper that does
// raw string -> tokens -> Command in one go.
func ParseLine(line string) (Command, error) {
	tokens, err := Lex(line)
	if err != nil {
		return nil, err
	}
	return ParseTokens(tokens)
}
