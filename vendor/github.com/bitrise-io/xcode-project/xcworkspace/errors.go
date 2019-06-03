package xcworkspace

import "fmt"

// SchemeNotFoundError represents that the given scheme was not found in the container
type SchemeNotFoundError struct {
	scheme    string
	container string
}

// Error implements the error interface
func (e SchemeNotFoundError) Error() string {
	return fmt.Sprintf("scheme %s not found in %s", e.scheme, e.container)
}

// IsSchemeNotFoundError reports whatever the given error is an instance of SchemeNotFoundError
func IsSchemeNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(SchemeNotFoundError)
	return ok
}
