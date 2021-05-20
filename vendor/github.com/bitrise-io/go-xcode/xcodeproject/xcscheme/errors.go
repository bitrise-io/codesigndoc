package xcscheme

import "fmt"

// NotFoundError represents that the given scheme was not found in the container
type NotFoundError struct {
	Scheme    string
	Container string
}

// Error implements the error interface
func (e NotFoundError) Error() string {
	return fmt.Sprintf("scheme %s not found in %s", e.Scheme, e.Container)
}

// IsNotFoundError reports whatever the given error is an instance of NotFoundError
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(NotFoundError)
	return ok
}
