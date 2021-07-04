package domain

import (
	"fmt"
)

// BasicDomainError is an abstraction to archive better error tracing
type BasicDomainError struct {
	op  string // domain operations stack for pretty error stacktrace
	err error  // the error itself

	parent *BasicDomainError
}

// Error stack up errors for better tracing
func Error(op string, err error) *BasicDomainError {
	if err == nil {
		return nil
	}

	domainError := BasicDomainError{op: op, err: err}
	if derr, ok := err.(*BasicDomainError); ok {
		derr.parent = &domainError
	}

	return &domainError
}

// Unwrap returns the err itself
func (e BasicDomainError) Unwrap() error {
	return e.err
}

// Error returns a string formatted error
func (e BasicDomainError) Error() string {
	return fmt.Sprintf("%s%s%s", e.op, e.separator(), e.err.Error())
}

func (e BasicDomainError) separator() string {
	if !e.isDomainErr() { // the root error that happened
		return ": "
	}
	return "->" // stacktrace
}

func (e BasicDomainError) isDomainErr() bool {
	if _, ok := e.err.(*BasicDomainError); ok {
		return true
	}
	return false
}
