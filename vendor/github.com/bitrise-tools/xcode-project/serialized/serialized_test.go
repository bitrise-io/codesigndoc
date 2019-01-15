package serialized

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeys(t *testing.T) {
	o := Object(map[string]interface{}{"key": "value", "key1": "value", "key2": "value"})
	keys := o.Keys()
	require.Equal(t, 3, len(keys))
	require.Contains(t, keys, "key")
	require.Contains(t, keys, "key1")
	require.Contains(t, keys, "key2")
}

func TestValue(t *testing.T) {
	o := Object(map[string]interface{}{"key": "value"})

	{
		v, err := o.Value("key")
		require.NoError(t, err)
		require.Equal(t, "value", v)
	}

	{
		v, err := o.Value("not_existing_key")
		require.EqualError(t, err, `key: string("not_existing_key") not found in: serialized.Object(serialized.Object{"key":"value"})`)
		require.Equal(t, nil, v)
	}
}

func TestString(t *testing.T) {
	o := Object(map[string]interface{}{"key": "value"})

	{
		v, err := o.String("key")
		require.NoError(t, err)
		require.Equal(t, "value", v)
	}

	{
		v, err := o.String("key2")
		require.EqualError(t, err, `key: string("key2") not found in: serialized.Object(serialized.Object{"key":"value"})`)
		require.Equal(t, "", v)
	}
}

func TestObject(t *testing.T) {
	o := Object(map[string]interface{}{"key": map[string]interface{}{"object_key": "object_value"}})

	{
		v, err := o.Object("key")
		require.NoError(t, err)
		require.Equal(t, Object(map[string]interface{}{"object_key": "object_value"}), v)
	}

	{
		v, err := o.Object("key2")
		require.EqualError(t, err, `key: string("key2") not found in: serialized.Object(serialized.Object{"key":map[string]interface {}{"object_key":"object_value"}})`)
		require.Equal(t, Object(nil), v)
	}
}

func TestStringSlice(t *testing.T) {
	o := Object{"buildConfigurations": []interface{}{"13E76E3B1F4AC90A0028096E", "13E76E3C1F4AC90A0028096E"}}

	{
		v, err := o.StringSlice("buildConfigurations")
		require.NoError(t, err)
		require.Equal(t, []string{"13E76E3B1F4AC90A0028096E", "13E76E3C1F4AC90A0028096E"}, v)
	}

	{
		v, err := o.StringSlice("key")
		require.EqualError(t, err, `key: string("key") not found in: serialized.Object(serialized.Object{"buildConfigurations":[]interface {}{"13E76E3B1F4AC90A0028096E", "13E76E3C1F4AC90A0028096E"}})`)
		require.Equal(t, []string(nil), v)
	}
}
