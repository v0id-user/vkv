package protocol

import (
    "bufio"
    "fmt"
)

// Encoder serializes Response objects into the wire protocol.
type Encoder struct {
    w *bufio.Writer
}

func NewEncoder(w *bufio.Writer) *Encoder {
    return &Encoder{w: w}
}

// Encode writes a single Response to the client.
func (e *Encoder) Encode(resp Response) error {
    switch r := resp.(type) {

    case RespOK:
        _, err := e.w.WriteString("OK\n")
        return err

    case RespValue:
        // VALUE <string>
        _, err := e.w.WriteString(fmt.Sprintf("VALUE %s\n", r.Value))
        return err

    case RespNil:
        _, err := e.w.WriteString("NIL\n")
        return err

    case RespErr:
        // ERR <msg>
        _, err := e.w.WriteString(fmt.Sprintf("ERR %s\n", r.Message))
        return err

    default:
        // Should never happen but just in case:
        _, err := e.w.WriteString("ERR internal encoder error\n")
        return err
    }
}
