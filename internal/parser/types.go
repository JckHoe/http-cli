package parser

import "net/http"

type CaptureRule struct {
	VariableName string
	JSONPath     string
}

type HTTPRequest struct {
	Method      string
	URL         string
	Headers     http.Header
	Body        string
	Name        string
	Description string
	LineNumber  int
	Variables   map[string]string
	Captures    []CaptureRule
}

type HTTPFile struct {
	Path     string
	Requests []HTTPRequest
	Variables map[string]string
}

type ParseError struct {
	Line    int
	Message string
}

func (e ParseError) Error() string {
	return e.Message
}