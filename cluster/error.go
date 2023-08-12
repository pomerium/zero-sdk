package cluster

import (
	"fmt"
)

// terminalError is an error that should not be retried
type terminalError struct {
	Err error
}

// Error implements error for terminalError
func (e *terminalError) Error() string {
	return fmt.Sprintf("terminal error: %v", e.Err)
}

// Unwrap implements errors.Unwrap for terminalError
func (e *terminalError) Unwrap() error {
	return e.Err
}

// IsTerminal implements TerminalError interface
// it may be used to check if an error is a terminal error in other packages
func (e *terminalError) IsTerminal() {}

// NewTerminalError creates a new terminal error that should not be retried
func NewTerminalError(err error) error {
	return &terminalError{Err: err}
}
