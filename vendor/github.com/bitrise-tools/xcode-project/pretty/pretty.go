package pretty

import (
	"encoding/json"
	"fmt"
)

// Object returns the json representation of the given parameter
//  if json marshal fails it returns the parameter's default format
func Object(o interface{}) string {
	b, err := json.MarshalIndent(o, "", "\t")
	if err != nil {
		return fmt.Sprint(o)
	}
	return string(b)
}
