package jobgroup

import "fmt"

// JobNotStartedError is returned when a job was never run.
//
// This usually indicates that the group's context ended before the job could start.
type JobNotStartedError struct {
	inner error
}

func newJobNotStartedError(inner error) *JobNotStartedError {
	return &JobNotStartedError{inner: inner}
}

// Error implements `error`.
func (e *JobNotStartedError) Error() string {
	return fmt.Sprintf("job could not be started: %v", e.inner)
}

// As implements the interface expected by `errors.As`.
func (e *JobNotStartedError) As(target any) bool {
	casted, ok := target.(**JobNotStartedError)
	if !ok {
		return false
	}

	*casted = e

	return true
}

// Unwrap implements the interface expected by `errors.Unwrap`.
func (e *JobNotStartedError) Unwrap() error {
	return e.inner
}
