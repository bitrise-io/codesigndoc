package serialized

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeyNotFoundError(t *testing.T) {
	err := NewKeyNotFoundError("key1", Object{"key": "value"})
	require.Equal(t, `key: string("key1") not found in: serialized.Object(serialized.Object{"key":"value"})`, err.Error())

	require.True(t, IsKeyNotFoundError(err))
	require.False(t, IsKeyNotFoundError(nil))

	require.False(t, IsTypeCastError(err))
}

func TestTypeCastError(t *testing.T) {
	err := NewTypeCastError("key1", "value", 0)
	require.Equal(t, `value: string("value") for key: string("key1") can not be casted to: int`, err.Error())

	require.True(t, IsTypeCastError(err))
	require.False(t, IsTypeCastError(nil))

	require.False(t, IsKeyNotFoundError(err))
}
