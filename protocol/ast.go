package protocol

// Command is the base interface for all commands.
// Each command type (Set, Get, Del) will implement this
type Command interface {
	Kind() string
}

// AST Node Types ---

type Set struct {
    Key   string
    Value string
}

func (s Set) Kind() string { return "SET" }


type Get struct {
	Key string
}

func (g Get) Kind() string { return "GET" }


type Del struct {
    Key string
}

func (d Del) Kind() string { return "DEL" }