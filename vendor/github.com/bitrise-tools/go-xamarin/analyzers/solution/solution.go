package solution

import (
	"bufio"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xamarin/analyzers/project"
	"github.com/bitrise-tools/go-xamarin/constants"
	"github.com/bitrise-tools/go-xamarin/utility"
)

const (
	solutionProjectsPattern = `Project\("{(?P<solution_id>[^"]*)}"\) = "(?P<project_name>[^"]*)", "(?P<project_path>[^"]*)", "{(?P<project_id>[^"]*)}"`

	solutionConfigurationPlatformsSectionStartPattern = `GlobalSection\(SolutionConfigurationPlatforms\) = preSolution`
	solutionConfigurationPlatformsSectionEndPattern   = `EndGlobalSection`
	solutionConfigurationPlatformPattern              = `(?P<config>[^|]*)\|(?P<platform>[^|]*) = (?P<m_config>[^|]*)\|(?P<m_platform>[^|]*)`

	projectConfigurationPlatformsSectionStartPattern = `GlobalSection\(ProjectConfigurationPlatforms\) = postSolution`
	projectConfigurationPlatformsSectionEndPattern   = `EndGlobalSection`
	projectConfigurationPlatformPattern              = `{(?P<project_id>.*)}.(?P<config>.*)\|(?P<platform>.*)\.Build.* = (?P<mapped_config>.*)\|(?P<mapped_platform>.*)`
)

// Model ...
type Model struct {
	Pth  string
	Name string
	ID   string

	ConfigMap map[string]string // Internal Configuartion|Platform - External Configuartion|Platform map

	ProjectMap map[string]project.Model // Project ID - Project Model map
}

// New ...
func New(pth string, loadProjects bool) (Model, error) {
	return analyzeSolution(pth, loadProjects)
}

// ConfigList ...
func (solution Model) ConfigList() []string {
	configList := []string{}
	for config := range solution.ConfigMap {
		configList = append(configList, config)
	}
	return configList
}

func analyzeSolution(pth string, analyzeProjects bool) (Model, error) {
	absPth, err := pathutil.AbsPath(pth)
	if err != nil {
		return Model{}, fmt.Errorf("Failed to expand path (%s), error: %s", pth, err)
	}

	fileName := filepath.Base(absPth)
	ext := filepath.Ext(absPth)
	fileName = strings.TrimSuffix(fileName, ext)

	solution := Model{
		Pth:        absPth,
		Name:       fileName,
		ConfigMap:  map[string]string{},
		ProjectMap: map[string]project.Model{},
	}

	isSolutionConfigurationPlatformsSection := false
	isProjectConfigurationPlatformsSection := false

	solutionDir := filepath.Dir(absPth)

	content, err := fileutil.ReadStringFromFile(absPth)
	if err != nil {
		return Model{}, fmt.Errorf("failed to read solution (%s), error: %s", absPth, err)
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Projects
		if matches := regexp.MustCompile(solutionProjectsPattern).FindStringSubmatch(line); len(matches) == 5 {
			ID := strings.ToUpper(matches[1])
			projectName := matches[2]
			projectID := strings.ToUpper(matches[4])
			projectRelativePth := utility.FixWindowsPath(matches[3])
			projectPth := filepath.Join(solutionDir, projectRelativePth)

			if strings.HasSuffix(projectPth, constants.CSProjExt) ||
				strings.HasSuffix(projectPth, constants.SHProjExt) ||
				strings.HasSuffix(projectPth, constants.FSProjExt) {

				project := project.Model{
					ID:   projectID,
					Name: projectName,
					Pth:  projectPth,

					ConfigMap: map[string]string{},
					Configs:   map[string]project.ConfigurationPlatformModel{},
				}
				solution.ProjectMap[projectID] = project
			}

			solution.ID = ID

			continue
		}

		// GlobalSection(SolutionConfigurationPlatforms) = preSolution
		if isSolutionConfigurationPlatformsSection {
			if match := regexp.MustCompile(solutionConfigurationPlatformsSectionEndPattern).FindString(line); match != "" {
				isSolutionConfigurationPlatformsSection = false
				continue
			}
		}

		if match := regexp.MustCompile(solutionConfigurationPlatformsSectionStartPattern).FindString(line); match != "" {
			isSolutionConfigurationPlatformsSection = true
			continue
		}

		if isSolutionConfigurationPlatformsSection {
			if matches := regexp.MustCompile(solutionConfigurationPlatformPattern).FindStringSubmatch(line); len(matches) == 5 {
				configuration := matches[1]
				platform := matches[2]

				mappedConfiguration := matches[3]
				mappedPlatform := matches[4]

				solution.ConfigMap[utility.ToConfig(configuration, platform)] = utility.ToConfig(mappedConfiguration, mappedPlatform)

				continue
			}
		}

		// GlobalSection(ProjectConfigurationPlatforms) = postSolution
		if isProjectConfigurationPlatformsSection {
			if match := regexp.MustCompile(projectConfigurationPlatformsSectionEndPattern).FindString(line); match != "" {
				isProjectConfigurationPlatformsSection = false
				continue
			}
		}

		if match := regexp.MustCompile(projectConfigurationPlatformsSectionStartPattern).FindString(line); match != "" {
			isProjectConfigurationPlatformsSection = true
			continue
		}

		if isProjectConfigurationPlatformsSection {
			if matches := regexp.MustCompile(projectConfigurationPlatformPattern).FindStringSubmatch(line); len(matches) == 6 {
				projectID := strings.ToUpper(matches[1])

				project, found := solution.ProjectMap[projectID]
				if !found {
					continue
				}

				solutionConfiguration := matches[2]
				solutionPlatform := matches[3]
				projectConfiguration := matches[4]
				projectPlatform := matches[5]
				if projectPlatform == "Any CPU" {
					projectPlatform = "AnyCPU"
				}

				project.ConfigMap[utility.ToConfig(solutionConfiguration, solutionPlatform)] = utility.ToConfig(projectConfiguration, projectPlatform)

				solution.ProjectMap[projectID] = project

				continue
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return Model{}, err
	}

	if analyzeProjects {
		projectMap := map[string]project.Model{}

		for projectID, proj := range solution.ProjectMap {
			projectDefinition, err := project.New(proj.Pth)
			if err != nil {
				return Model{}, fmt.Errorf("failed to analyze project (%s), error: %s", proj.Pth, err)
			}

			projectDefinition.Name = proj.Name
			projectDefinition.Pth = proj.Pth
			projectDefinition.ConfigMap = proj.ConfigMap

			projectMap[projectID] = projectDefinition
		}

		solution.ProjectMap = projectMap
	}

	return solution, nil
}
