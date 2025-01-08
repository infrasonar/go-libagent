package libagent

import (
	"errors"
	"fmt"
)

type Severity string

const (
	Low    Severity = "LOW"
	Medium Severity = "MEDIUM"
	High   Severity = "HIGH"
)

type CheckError struct {
	Sev Severity
	Err error
}

func (r *CheckError) Error() string {
	return r.Err.Error()
}

func ErrLow(text string) *CheckError {
	return &CheckError{
		Sev: Low,
		Err: errors.New(text),
	}
}

func ErrLowf(format string, a ...any) *CheckError {
	return &CheckError{
		Sev: Low,
		Err: fmt.Errorf(format, a...),
	}
}

func ErrMedium(text string) *CheckError {
	return &CheckError{
		Sev: Medium,
		Err: errors.New(text),
	}
}

func ErrMediumf(format string, a ...any) *CheckError {
	return &CheckError{
		Sev: Medium,
		Err: fmt.Errorf(format, a...),
	}
}

func ErrHigh(text string) *CheckError {
	return &CheckError{
		Sev: High,
		Err: errors.New(text),
	}
}

func ErrHighf(format string, a ...any) *CheckError {
	return &CheckError{
		Sev: High,
		Err: fmt.Errorf(format, a...),
	}
}
