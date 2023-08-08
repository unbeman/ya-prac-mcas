package metrics

import "errors"

var (
	ErrInvalidType  = errors.New("invalid type")
	ErrInvalidValue = errors.New("invalid value")
	ErrParseJSON    = errors.New("can't parse")
	ErrParseURI     = errors.New("can't parse")
)
