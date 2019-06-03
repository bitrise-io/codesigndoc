package utility

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/go-utils/sliceutil"
)

// CreateIconDescriptors returns a sorted array of generated unique icon names, and the original file path.
func CreateIconDescriptors(absoluteIconPaths []string, basepath string) (models.Icons, error) {
	absoluteIconPaths = sliceutil.UniqueStringSlice(absoluteIconPaths)
	sort.Strings(absoluteIconPaths)

	icons := models.Icons{}
	for _, iconPath := range absoluteIconPaths {
		relativePath, err := filepath.Rel(basepath, iconPath)
		if err != nil {
			return nil, err
		}
		hash := sha256.Sum256([]byte(relativePath))
		hashStr := fmt.Sprintf("%x", hash) + filepath.Ext(iconPath)

		icons = append(icons, models.Icon{
			Filename: hashStr,
			Path:     iconPath,
		})
	}
	return icons, nil
}
