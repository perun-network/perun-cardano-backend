package types

import "fmt"

type DecodeEventError struct {
	Tag         string
	ExpectedLen int
	ActualLen   int
}

func NewDecodeEventError(tag string, expectedLen int, actualLen int) *DecodeEventError {
	return &DecodeEventError{
		Tag:         tag,
		ExpectedLen: expectedLen,
		ActualLen:   actualLen,
	}
}

func (e DecodeEventError) Error() string {
	return fmt.Sprintf("invalid amount of ChannelDatums received in %s event. Expected: %d, Actual: %d",
		e.Tag,
		e.ExpectedLen,
		e.ActualLen)
}
