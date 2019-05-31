package ios

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type assetIcon struct {
	Size     string
	Filename string
}

type assetInfo struct {
	Version int
	Author  string
}
type assetMetadata struct {
	Images []assetIcon
	Info   assetInfo
}

type appIcon struct {
	Size     int
	Filename string
}

func parseResourceSet(resourceSetPath string) (appIcon, bool, error) {
	const resourceMetadataFileName = "Contents.json"
	file, err := os.Open(filepath.Join(resourceSetPath, resourceMetadataFileName))
	if err != nil {
		return appIcon{}, false, fmt.Errorf("failed to open file, error: %s", err)
	}

	appIcons, err := parseResourceSetMetadata(file)
	if err != nil {
		return appIcon{}, false, fmt.Errorf("failed to parse asset metadata, error: %s", err)
	}

	if len(appIcons) == 0 {
		return appIcon{}, false, nil
	}
	largestIcon := appIcons[0]
	for _, icon := range appIcons {
		if icon.Size > largestIcon.Size {
			largestIcon = icon
		}
	}
	return largestIcon, true, nil
}

func parseResourceSetMetadata(input io.Reader) ([]appIcon, error) {
	var decoded assetMetadata
	if err := json.NewDecoder(input).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("failed to decode asset metadata file, error: %s", err)
	}

	if decoded.Info.Version != 1 {
		return nil, fmt.Errorf("unsupported asset metadata version")
	}
	var icons []appIcon
	for _, icon := range decoded.Images {
		sizeParts := strings.Split(icon.Size, "x")
		if len(sizeParts) != 2 {
			return nil, fmt.Errorf("invalid image size format")
		}
		iconSize, err := strconv.ParseFloat(sizeParts[0], 32)
		if err != nil {
			return nil, fmt.Errorf("invalid image size, error: %s", err)
		}
		// If icon is not set for a given usage, filaname key is missing
		if icon.Filename != "" {
			icons = append(icons, appIcon{
				Size:     int(iconSize),
				Filename: icon.Filename,
			})
		}
	}
	return icons, nil
}
