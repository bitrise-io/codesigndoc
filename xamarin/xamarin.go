package xamarin

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/progress"
	"github.com/bitrise-tools/go-xamarin/constants"
)

// CommandModel ...
type CommandModel struct {
	SolutionFilePath string
	ProjectName      string
	Configuration    string
	Platform         string
}

// SetConfigurationPlatformCombination - `configPlatformCombination` should be a composite string
// with the format: "CONFIGURATION|PLATFORM"
// e.g.: Release|iPhone
func (xamarinCmd *CommandModel) SetConfigurationPlatformCombination(configPlatformCombination string) error {
	split := strings.Split(configPlatformCombination, "|")
	if len(split) != 2 {
		return fmt.Errorf("invalid configuration-platform combination (%s), should include exactly one pipe (|) character", configPlatformCombination)
	}
	xamarinCmd.Configuration = split[0]
	xamarinCmd.Platform = split[1]
	return nil
}

// GenerateArchive ...
func (xamarinCmd CommandModel) GenerateArchive() (string, string, error) {
	cmdOut := ""
	archivePth := ""
	var err error

	progress.SimpleProgress(".", 1*time.Second, func() {
		archivePth, cmdOut, err = xamarinCmd.RunBuildCommand()
	})
	fmt.Println()

	if err != nil {
		return "", cmdOut, fmt.Errorf("Failed to Archive, error: %s", err)
	}

	return archivePth, cmdOut, nil
}

// RunBuildCommand ...
func (xamarinCmd CommandModel) RunBuildCommand() (string, string, error) {
	// STO: https://stackoverflow.com/a/19534376/5842489
	// if your project has a . (dot) in its name, replace it with a _ (underscore) when specifying it with /t
	projectName := strings.Replace(xamarinCmd.ProjectName, ".", "_", -1)

	archivesBeforeBuild, err := listArchives()
	if err != nil {
		return "", "", fmt.Errorf("failed to list archives, error: %s", err)
	}

	cmdArgs := []string{constants.MsbuildPath,
		xamarinCmd.SolutionFilePath,
		fmt.Sprintf("/p:Configuration=%s", xamarinCmd.Configuration),
		fmt.Sprintf("/p:Platform=%s", xamarinCmd.Platform),
		fmt.Sprintf("/p:ArchiveOnBuild=true"),
		fmt.Sprintf("/t:%s", projectName),
	}

	log.Infof("$ %s", command.PrintableCommandArgs(true, cmdArgs))
	cmd, err := command.NewFromSlice(cmdArgs)
	if err != nil {
		return "", "", fmt.Errorf("Failed to create Xamarin command, error: %s", err)
	}
	xamarinBuildOutput, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return "", xamarinBuildOutput, fmt.Errorf("Failed to run Xamarin command, error: %s", err)
	}

	log.Debugf("xamarinBuildOutput: %s", xamarinBuildOutput)

	archivesAfterBuild, err := listArchives()
	if err != nil {
		return "", "", fmt.Errorf("failed to list archives, error: %s", err)
	}

	newArchives := []string{}
	for _, archiveAfterBuild := range archivesAfterBuild {
		isNew := true
		for _, archiveBeforeBuild := range archivesBeforeBuild {
			if archiveAfterBuild == archiveBeforeBuild {
				isNew = false
				break
			}
		}
		if isNew {
			newArchives = append(newArchives, archiveAfterBuild)
		}
	}

	if len(newArchives) == 0 {
		return "", xamarinBuildOutput, errors.New("No archive generated during the build")
	} else if len(newArchives) > 1 {
		return "", xamarinBuildOutput, errors.New("multiple archives generated during the build")
	}

	return newArchives[0], xamarinBuildOutput, nil
}

func listArchives() ([]string, error) {
	userHomeDir := os.Getenv("HOME")
	if userHomeDir == "" {
		return []string{}, errors.New("failed to get user home dir")
	}
	xcodeArchivesDir := filepath.Join(userHomeDir, "Library/Developer/Xcode/Archives")
	if exist, err := pathutil.IsDirExists(xcodeArchivesDir); err != nil {
		return []string{}, err
	} else if !exist {
		return []string{}, fmt.Errorf("no default Xcode archive path found at: %s", xcodeArchivesDir)
	}

	archives := []string{}
	if walkErr := filepath.Walk(xcodeArchivesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".xcarchive" {
			archives = append(archives, path)
		}

		return nil
	}); walkErr != nil {
		return []string{}, fmt.Errorf("failed to find archives, error: %s", walkErr)
	}

	return archives, nil
}
