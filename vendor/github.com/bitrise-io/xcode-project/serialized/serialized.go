package serialized

// Object ...
type Object map[string]interface{}

// Keys ...
func (o Object) Keys() []string {
	var keys []string
	for key := range o {
		keys = append(keys, key)
	}
	return keys
}

// Value ...
func (o Object) Value(key string) (interface{}, error) {
	value, ok := o[key]
	if !ok {
		return nil, NewKeyNotFoundError(key, o)
	}
	return value, nil
}

// String ...
func (o Object) String(key string) (string, error) {
	value, err := o.Value(key)
	if err != nil {
		return "", err
	}

	casted, ok := value.(string)
	if !ok {
		return "", NewTypeCastError(key, value, "")
	}

	return casted, nil
}

// Int64 returns a value with int64 type from the map
func (o Object) Int64(key string) (int64, error) {
	value, err := o.Value(key)
	if err != nil {
		return -1, err
	}

	casted, ok := value.(int64)
	if !ok {
		return -1, NewTypeCastError(key, value, 0)
	}

	return casted, nil
}

// StringSlice ...
func (o Object) StringSlice(key string) ([]string, error) {
	value, err := o.Value(key)
	if err != nil {
		return nil, err
	}

	casted, ok := value.([]interface{})
	if !ok {
		return nil, NewTypeCastError(key, value, []interface{}{})
	}

	slice := []string{}
	for _, v := range casted {
		item, ok := v.(string)
		if !ok {
			return nil, NewTypeCastError(key, casted, "")
		}

		slice = append(slice, item)
	}

	return slice, nil
}

// Object ...
func (o Object) Object(key string) (Object, error) {
	value, err := o.Value(key)
	if err != nil {
		return nil, err
	}

	casted, ok := value.(map[string]interface{})
	if !ok {
		return nil, NewTypeCastError(key, value, map[string]interface{}{})
	}

	return casted, nil
}
