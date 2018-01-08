package xamarin

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/progress"
	"github.com/bitrise-tools/go-xamarin/builder"
	"github.com/bitrise-tools/go-xamarin/constants"
	"github.com/bitrise-tools/go-xamarin/tools/buildtools"
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
		return "", cmdOut, err
	}

	return archivePth, cmdOut, nil
}

// RunBuildCommand ...
func (xamarinCmd CommandModel) RunBuildCommand() (string, string, error) {
	xamarinBuilder, err := builder.New(xamarinCmd.SolutionFilePath, []constants.SDK{constants.SDKIOS}, buildtools.Msbuild)
	if err != nil {
		return "", "", err
	}

	var outWriter bytes.Buffer
	xamarinBuilder.SetOutputs(&outWriter, &outWriter)

	archivesBeforeBuild, err := listArchives()
	if err != nil {
		return "", "", fmt.Errorf("failed to list before build archives, error: %s", err)
	}

	callback := func(_ string, projectName string, _ constants.SDK, _ constants.TestFramework, commandStr string, alreadyPerformed bool) {
		log.Printf("\nBuilding project: %s", projectName)
		log.Infof("$ %s", commandStr)
		if alreadyPerformed {
			log.Warnf("build command already performed, skipping...")
		}
	}

	warnings, err := xamarinBuilder.BuildAllProjects(xamarinCmd.Configuration, xamarinCmd.Platform, false, nil, callback)
	xamarinBuildOutput := outWriter.String()

	log.Debugf("xamarinBuildOutput: %s", xamarinBuildOutput)

	if len(warnings) > 0 {
		log.Warnf("Build warnings:")
		for _, warning := range warnings {
			log.Warnf(warning)
		}
	}
	if err != nil {
		return "", xamarinBuildOutput, err
	}

	archivesAfterBuild, err := listArchives()
	if err != nil {
		return "", "", fmt.Errorf("failed to list after build archives, error: %s", err)
	}

	var archivesDuringBuild []string
	for _, afterArchive := range archivesAfterBuild {
		generatedDuringBuild := true
		for _, beforeArchive := range archivesBeforeBuild {
			if beforeArchive == afterArchive {
				generatedDuringBuild = false
				break
			}
		}
		if generatedDuringBuild {
			archivesDuringBuild = append(archivesDuringBuild, afterArchive)
		}
	}

	if len(archivesDuringBuild) == 0 {
		return "", xamarinBuildOutput, fmt.Errorf("failed to find the xcarchive generated during the build")
	} else if len(archivesDuringBuild) > 1 {
		return "", xamarinBuildOutput, fmt.Errorf("multiple xcarchives generated during the build")
	}

	return archivesDuringBuild[0], xamarinBuildOutput, nil
}

func listArchives() ([]string, error) {
	userHomeDir, found := os.LookupEnv("HOME")
	if !found {
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
