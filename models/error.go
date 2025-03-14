package models

import (
	"encoding/json"
	"fmt"
)

type ErrorKind int

const (
	ErrKindInternal ErrorKind = iota
	ErrKindInvalid
	ErrKindMissing
)

func (e ErrorKind) String() string {
	switch e {
	case ErrKindInvalid:
		return "invalid"
	case ErrKindMissing:
		return "missing"
	case ErrKindInternal:
		fallthrough
	default:
		return "internal"
	}
}

func (e ErrorKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

type Error struct {
	Kind    ErrorKind `json:"kind"`
	Message string    `json:"message"`
	Cause   error     `json:"-"`
}

func (e Error) Error() string {
	return e.Message
}

func ErrWithCause(kind ErrorKind, message string, cause error) Error {
	return Error{Kind: kind, Message: message, Cause: cause}
}

func ErrInvalid(message string) Error {
	return ErrWithCause(ErrKindInvalid, message, nil)
}

func ErrInvalidf(message string, args ...any) Error {
	return ErrWithCause(ErrKindInvalid, fmt.Sprintf(message, args...), nil)
}

func ErrInvalidWithCause(message string, cause error) Error {
	return ErrWithCause(ErrKindInvalid, message, cause)
}

func ErrMissing(message string) Error {
	return ErrWithCause(ErrKindMissing, message, nil)
}

func ErrMissingf(message string, args ...any) Error {
	return ErrWithCause(ErrKindMissing, fmt.Sprintf(message, args...), nil)
}

func ErrInternalWithCause(message string, cause error) Error {
	return ErrWithCause(ErrKindInternal, message, cause)
}
