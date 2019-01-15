package pretty

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestObject(t *testing.T) {
	o := map[string]interface{}{
		"key": "value",
		"slice_key": []string{
			"item1",
			"item2",
		},
	}
	require.Equal(t, `{
	"key": "value",
	"slice_key": [
		"item1",
		"item2"
	]
}`, Object(o), Object(o))
}
