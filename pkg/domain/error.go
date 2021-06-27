package domain

import (
	"encoding/json"
	"fmt"
)

type Params map[string]interface{}

type BasicDomainError struct {
	op     string  // domain operations stack for pretty error stacktrace
	err    error   // the error itself
	params *Params // params sent to the last called function which triggered the error (optional)

	parent *BasicDomainError
}

func Error(op string, err error, opts ...Params) *BasicDomainError {
	if err == nil {
		return nil
	}

	domainError := BasicDomainError{op: op, err: err}
	if len(opts) > 0 {
		domainError.params = &opts[0]
	}
	if derr, ok := err.(*BasicDomainError); ok {
		derr.parent = &domainError
	}

	return &domainError
}

func (e BasicDomainError) Unwrap() error {
	return e.err
}

func (e BasicDomainError) Error() string {
	return fmt.Sprintf("%s%s%s%s", e.op, e.calledWithParams(), e.separator(), e.err.Error())
}

func (e BasicDomainError) calledWithParams() string {
	if !e.hasParentErr() || e.parent.params == nil {
		return "()"
	}

	keyValues, _ := json.Marshal(e.parent.params)
	keyValues = keyValues[1 : len(keyValues)-1] // drop json braces
	return fmt.Sprintf("(%s)", keyValues)
}

func (e BasicDomainError) hasParentErr() bool {
	return e.parent != nil
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
