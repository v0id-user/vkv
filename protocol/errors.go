package protocol

import "fmt"

// ProtocolError represents an error caused by invalid client input.
// Examples: unknown command, wrong argument count, bad characters, etc.
type ProtocolError struct {
    Msg string
}

func (e ProtocolError) Error() string {
    return e.Msg
}

// NewProtocolError creates a protocol-level error.
func NewProtocolError(msg string) error {
    return ProtocolError{Msg: msg}
}

// IsProtocolError tells the server whether this error
// came from the client or the system.
func IsProtocolError(err error) bool {
    _, ok := err.(ProtocolError)
    return ok
}

// Predeclared common protocol errors
var (
    ErrEmptyCommand = ProtocolError{"empty command"}
    ErrUnknownCommand = ProtocolError{"unknown command"}
    ErrBadArgumentCount = ProtocolError{"wrong number of arguments"}
)

// Helper to wrap dynamic messages
func Errorf(format string, args ...any) error {
    return ProtocolError{Msg: fmt.Sprintf(format, args...)}
}