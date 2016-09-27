package solution

import (
	"bufio"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-tools/go-xamarin/project"
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
	ID  string
	Pth string

	ConfigMap map[string]string

	ProjectMap map[string]project.Model
}

// New ...
func New(pth string, loadProjects bool) (Model, error) {
	return analyzeSolution(pth, loadProjects)
}

func (solution Model) String() string {
	s := ""
	s += fmt.Sprintf("ID: %s\n", solution.ID)
	s += fmt.Sprintf("Pth: %s\n", solution.Pth)
	s += "\n"
	s += fmt.Sprintf("ConfigMap:\n")
	s += fmt.Sprintf("%v\n", solution.ConfigMap)
	s += "\n"
	s += fmt.Sprintf("ProjectMap:\n")
	s += fmt.Sprintf("%v\n", solution.ProjectMap)
	s += "\n"
	return s
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
	solution := Model{
		Pth:        pth,
		ConfigMap:  map[string]string{},
		ProjectMap: map[string]project.Model{},
	}

	isSolutionConfigurationPlatformsSection := false
	isProjectConfigurationPlatformsSection := false

	solutionDir := filepath.Dir(pth)

	content, err := fileutil.ReadStringFromFile(pth)
	if err != nil {
		return Model{}, fmt.Errorf("failed to read solution (%s), error: %s", pth, err)
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Projects
		if matches := regexp.MustCompile(solutionProjectsPattern).FindStringSubmatch(line); len(matches) == 5 {
			ID := matches[1]
			projectName := matches[2]
			projectID := matches[4]
			projectRelativePth := utility.FixWindowsPath(matches[3])
			projectPth := filepath.Join(solutionDir, projectRelativePth)

			solution.ID = ID

			project := project.Model{
				ID:   projectID,
				Name: projectName,
				Pth:  projectPth,

				ConfigMap: map[string]string{},
				Configs:   map[string]project.ConfigurationPlatformModel{},
			}
			solution.ProjectMap[projectID] = project

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
				projectID := matches[1]
				solutionConfiguration := matches[2]
				solutionPlatform := matches[3]
				projectConfiguration := matches[4]
				projectPlatform := matches[5]
				if projectPlatform == "Any CPU" {
					projectPlatform = "AnyCPU"
				}

				project, found := solution.ProjectMap[projectID]
				if !found {
					return Model{}, fmt.Errorf("no project found with ID: %s", projectID)
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
