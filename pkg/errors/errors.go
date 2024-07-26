package errors

import "errors"

func newError(msg string) error {
	return errors.New(msg)
}

var (
	ErrUnknownAirCompany  = newError("Unknown Air Company")
	ErrAirCompanyNotFound = newError("Air Company not found")
)
