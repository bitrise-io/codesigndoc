package utility

import (
	"fmt"
	"strings"
)

// ToConfig ...
func ToConfig(configuration, platform string) string {
	return fmt.Sprintf("%s|%s", configuration, platform)
}

// FixWindowsPath ...
func FixWindowsPath(pth string) string {
	return strings.Replace(pth, `\`, "/", -1)
}

// SplitAndStripList ...
func SplitAndStripList(list, separator string) []string {
	split := strings.Split(list, separator)
	elements := []string{}
	for _, s := range split {
		elements = append(elements, strings.TrimSpace(s))
	}
	return elements
}
