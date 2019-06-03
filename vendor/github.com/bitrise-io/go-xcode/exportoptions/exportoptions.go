package exportoptions

import (
	"fmt"
	"path/filepath"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"howett.net/plist"
)

// ExportOptions ...
type ExportOptions interface {
	Hash() map[string]interface{}
	String() (string, error)
	WriteToFile(pth string) error
	WriteToTmpFile() (string, error)
}

// WritePlistToFile ...
func WritePlistToFile(options map[string]interface{}, pth string) error {
	plistBytes, err := plist.MarshalIndent(options, plist.XMLFormat, "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal export options model, error: %s", err)
	}
	if err := fileutil.WriteBytesToFile(pth, plistBytes); err != nil {
		return fmt.Errorf("failed to write export options, error: %s", err)
	}

	return nil
}

// WritePlistToTmpFile ...
func WritePlistToTmpFile(options map[string]interface{}) (string, error) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("output")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir, error: %s", err)
	}
	pth := filepath.Join(tmpDir, "exportOptions.plist")

	if err := WritePlistToFile(options, pth); err != nil {
		return "", fmt.Errorf("failed to write to file options, error: %s", err)
	}

	return pth, nil
}
