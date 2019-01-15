package serialized

import "fmt"

// KeyNotFoundError ...
type KeyNotFoundError struct {
	key    string
	object Object
}

// Error ...
func (e KeyNotFoundError) Error() string {
	return fmt.Sprintf("key: %T(%#v) not found in: %T(%#v)", e.key, e.key, e.object, e.object)
}

// NewKeyNotFoundError ...
func NewKeyNotFoundError(key string, object Object) KeyNotFoundError {
	return KeyNotFoundError{key: key, object: object}
}

// IsKeyNotFoundError ...
func IsKeyNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(KeyNotFoundError)
	return ok
}

// TypeCastError ...
type TypeCastError struct {
	key          string
	value        interface{}
	expectedType string
}

// NewTypeCastError ...
func NewTypeCastError(key string, value interface{}, expected interface{}) TypeCastError {
	return TypeCastError{key: key, value: value, expectedType: fmt.Sprintf("%T", expected)}
}

// IsTypeCastError ...
func IsTypeCastError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(TypeCastError)
	return ok
}

// Error ...
func (e TypeCastError) Error() string {
	return fmt.Sprintf("value: %T(%#v) for key: %T(%#v) can not be casted to: %s", e.value, e.value, e.key, e.key, e.expectedType)
}
