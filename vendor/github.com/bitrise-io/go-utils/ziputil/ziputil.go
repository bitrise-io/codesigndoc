package ziputil

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
)

// ZipDir ...
func ZipDir(sourceDirPth, destinationZipPth string, isContentOnly bool) error {
	if exist, err := pathutil.IsDirExists(sourceDirPth); err != nil {
		return err
	} else if !exist {
		return fmt.Errorf("dir (%s) not exist", sourceDirPth)
	}

	workDir := filepath.Dir(sourceDirPth)
	zipTarget := filepath.Base(sourceDirPth)

	if isContentOnly {
		workDir = sourceDirPth
		zipTarget = "."
	}

	return internalZipDir(destinationZipPth, zipTarget, workDir)

}

// ZipDirs ...
func ZipDirs(sourceDirPths []string, destinationZipPth string) error {
	for _, path := range sourceDirPths {
		if exist, err := pathutil.IsDirExists(path); err != nil {
			return err
		} else if !exist {
			return fmt.Errorf("directory (%s) not exist", path)
		}
	}

	tempDir, err := pathutil.NormalizedOSTempDirPath("zip")
	if err != nil {
		return err
	}

	defer func() {
		if err = os.RemoveAll(tempDir); err != nil {
			log.Fatal(err)
		}
	}()

	for _, path := range sourceDirPths {
		err := command.CopyDir(path, tempDir, false)
		if err != nil {
			return err
		}
	}

	return internalZipDir(destinationZipPth, ".", tempDir)
}

func internalZipDir(destinationZipPth, zipTarget, workDir string) error {
	// -r - Travel the directory structure recursively
	// -T - Test the integrity of the new zip file
	// -y - Store symbolic links as such in the zip archive, instead of compressing and storing the file referred to by the link
	cmd := command.New("/usr/bin/zip", "-rTy", destinationZipPth, zipTarget)
	cmd.SetDir(workDir)
	if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		return fmt.Errorf("command: (%s) failed, output: %s, error: %s", cmd.PrintableCommandArgs(), out, err)
	}

	return nil
}

// ZipFile ...
func ZipFile(sourceFilePth, destinationZipPth string) error {
	return ZipFiles([]string{sourceFilePth}, destinationZipPth)
}

// ZipFiles ...
func ZipFiles(sourceFilePths []string, destinationZipPth string) error {
	for _, path := range sourceFilePths {
		if exist, err := pathutil.IsPathExists(path); err != nil {
			return err
		} else if !exist {
			return fmt.Errorf("file (%s) not exist", path)
		}
	}

	// -T - Test the integrity of the new zip file
	// -y - Store symbolic links as such in the zip archive, instead of compressing and storing the file referred to by the link
	// -j - Do not recreate the directory structure inside the zip. Kind of equivalent of copying all the files in one folder and zipping it.
	parameters := []string{"-Tyj", destinationZipPth}
	parameters = append(parameters, sourceFilePths...)
	cmd := command.New("/usr/bin/zip", parameters...)
	if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		return fmt.Errorf("command: (%s) failed, output: %s, error: %s", cmd.PrintableCommandArgs(), out, err)
	}

	return nil
}

// UnZip ...
func UnZip(zip, intoDir string) error {
	cmd := command.New("/usr/bin/unzip", zip, "-d", intoDir)
	if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		return fmt.Errorf("command: (%s) failed, output: %s, error: %s", cmd.PrintableCommandArgs(), out, err)
	}

	return nil
}
