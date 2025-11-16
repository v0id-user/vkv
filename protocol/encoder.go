package protocol

// Response is the output of the engine.
// Encoder serializes this for the client.
type Response interface {
	Kind() string
}

// OK
type RespOK struct{}

func (RespOK) Kind() string { return "OK" }

func ResponseOK() Response {
	return RespOK{}
}

// Value
type RespValue struct {
	Value string
}

func (RespValue) Kind() string { return "VALUE" }

func ResponseValue(val string) Response {
	return RespValue{Value: val}
}

// NIL (key not found)
type RespNil struct{}

func (RespNil) Kind() string { return "NIL" }

func ResponseNil() Response {
	return RespNil{}
}

// Error
type RespErr struct {
	Message string
}

func (RespErr) Kind() string { return "ERR" }

func ResponseErr(msg string) Response {
	return RespErr{Message: msg}
}
