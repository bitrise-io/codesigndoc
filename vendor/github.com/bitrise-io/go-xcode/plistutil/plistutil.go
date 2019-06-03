package plistutil

import (
	"errors"
	"time"

	"github.com/bitrise-io/go-utils/fileutil"
	"howett.net/plist"
)

// PlistData ...
type PlistData map[string]interface{}

// NewPlistDataFromContent ...
func NewPlistDataFromContent(plistContent string) (PlistData, error) {
	var data PlistData
	if _, err := plist.Unmarshal([]byte(plistContent), &data); err != nil {
		return PlistData{}, err
	}
	return data, nil
}

// NewPlistDataFromFile ...
func NewPlistDataFromFile(plistPth string) (PlistData, error) {
	content, err := fileutil.ReadStringFromFile(plistPth)
	if err != nil {
		return PlistData{}, err
	}
	return NewPlistDataFromContent(content)
}

// GetString ...
func (data PlistData) GetString(forKey string) (string, bool) {
	value, ok := data[forKey]
	if !ok {
		return "", false
	}

	casted, ok := value.(string)
	if !ok {
		return "", false
	}

	return casted, true
}

// GetUInt64 ...
func (data PlistData) GetUInt64(forKey string) (uint64, bool) {
	value, ok := data[forKey]
	if !ok {
		return 0, false
	}

	casted, ok := value.(uint64)
	if !ok {
		return 0, false
	}
	return casted, true
}

// GetFloat64 ...
func (data PlistData) GetFloat64(forKey string) (float64, bool) {
	value, ok := data[forKey]
	if !ok {
		return 0, false
	}

	casted, ok := value.(float64)
	if !ok {
		return 0, false
	}
	return casted, true
}

// GetBool ...
func (data PlistData) GetBool(forKey string) (bool, bool) {
	value, ok := data[forKey]
	if !ok {
		return false, false
	}

	casted, ok := value.(bool)
	if !ok {
		return false, false
	}

	return casted, true
}

// GetTime ...
func (data PlistData) GetTime(forKey string) (time.Time, bool) {
	value, ok := data[forKey]
	if !ok {
		return time.Time{}, false
	}

	casted, ok := value.(time.Time)
	if !ok {
		return time.Time{}, false
	}
	return casted, true
}

// GetUInt64Array ...
func (data PlistData) GetUInt64Array(forKey string) ([]uint64, bool) {
	value, ok := data[forKey]
	if !ok {
		return nil, false
	}

	if casted, ok := value.([]uint64); ok {
		return casted, true
	}

	casted, ok := value.([]interface{})
	if !ok {
		return nil, false
	}

	array := []uint64{}
	for _, v := range casted {
		casted, ok := v.(uint64)
		if !ok {
			return nil, false
		}

		array = append(array, casted)
	}
	return array, true
}

// GetStringArray ...
func (data PlistData) GetStringArray(forKey string) ([]string, bool) {
	value, ok := data[forKey]
	if !ok {
		return nil, false
	}

	if casted, ok := value.([]string); ok {
		return casted, true
	}

	casted, ok := value.([]interface{})
	if !ok {
		return nil, false
	}

	array := []string{}
	for _, v := range casted {
		casted, ok := v.(string)
		if !ok {
			return nil, false
		}

		array = append(array, casted)
	}
	return array, true
}

// GetByteArrayArray ...
func (data PlistData) GetByteArrayArray(forKey string) ([][]byte, bool) {
	value, ok := data[forKey]
	if !ok {
		return nil, false
	}

	if casted, ok := value.([][]byte); ok {
		return casted, true
	}

	casted, ok := value.([]interface{})
	if !ok {
		return nil, false
	}

	array := [][]byte{}
	for _, v := range casted {
		casted, ok := v.([]byte)
		if !ok {
			return nil, false
		}

		array = append(array, casted)
	}
	return array, true
}

// GetMapStringInterface ...
func (data PlistData) GetMapStringInterface(forKey string) (PlistData, bool) {
	value, ok := data[forKey]
	if !ok {
		return nil, false
	}

	if casted, ok := value.(map[string]interface{}); ok {
		return casted, true
	}
	return nil, false
}

func castToMapStringInterfaceArray(obj interface{}) ([]PlistData, error) {
	array, ok := obj.([]interface{})
	if !ok {
		return nil, errors.New("failed to cast to []interface{}")
	}

	var casted []PlistData
	for _, item := range array {
		mapStringInterface, ok := item.(map[string]interface{})
		if !ok {
			return nil, errors.New("failed to cast to map[string]interface{}")
		}
		casted = append(casted, mapStringInterface)
	}

	return casted, nil
}

// GetMapStringInterfaceArray ...
func (data PlistData) GetMapStringInterfaceArray(forKey string) ([]PlistData, bool) {
	value, ok := data[forKey]
	if !ok {
		return nil, false
	}
	mapStringInterfaceArray, err := castToMapStringInterfaceArray(value)
	if err != nil {
		return nil, false
	}
	return mapStringInterfaceArray, true
}
