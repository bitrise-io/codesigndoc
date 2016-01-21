package regexputil

import (
	"errors"
	"regexp"
)

// NamedFindStringSubmatch ...
func NamedFindStringSubmatch(rexp *regexp.Regexp, text string) (map[string]string, error) {
	match := rexp.FindStringSubmatch(text)
	if match == nil {
		return map[string]string{}, errors.New("No match found")
	}
	result := map[string]string{}
	for i, name := range rexp.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}
	return result, nil
}
