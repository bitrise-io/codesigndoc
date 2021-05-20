package ios

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"encoding/json"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-xcode/xcodeproj"
)

const podfileBase = "Podfile"

// AllowPodfileBaseFilter ...
var AllowPodfileBaseFilter = pathutil.BaseFilter(podfileBase, true)

type podfileParser struct {
	podfilePth                string
	suppressPodFileParseError bool
}

func (podfileParser podfileParser) getTargetDefinitionProjectMap(cocoapodsVersion string) (map[string]string, error) {
	gemfileCocoapodsVersion := ""
	if cocoapodsVersion != "" {
		gemfileCocoapodsVersion = fmt.Sprintf(`, '%s'`, cocoapodsVersion)
	}

	gemfileContent := fmt.Sprintf(`source 'https://rubygems.org'
gem 'cocoapods'%s
gem 'cocoapods-core'
gem 'json'
`, gemfileCocoapodsVersion)

	// returns target - project map, if xcodeproj defined in the Podfile
	// return empty string if no xcodeproj defined in the Podfile
	rubyScriptContent := `require 'cocoapods-core'
require 'json'

begin
	podfile_path = ENV['PODFILE_PATH']
	# In case of relative require in the Podfile, change working directory
	# For example: require_relative '../node_modules/react-native/scripts/react_native_pods'
	Dir.chdir(File.dirname(podfile_path))
	podfile = Pod::Podfile.from_file(podfile_path)
	targets = podfile.target_definitions
	
	puts "#{{}.to_json}" unless targets

	target_project_map = {}
	targets.each do |name, target_definition|
		next unless target_definition.user_project_path
		target_project_map[name] = target_definition.user_project_path
	end

	puts "#{{ :data => target_project_map }.to_json}"
rescue => e
	puts "#{{ :error => "#{e.to_s} Reason: #{e.message}"}.to_json}"
end
`

	absPodfilePth, err := filepath.Abs(podfileParser.podfilePth)
	if err != nil {
		return map[string]string{}, fmt.Errorf("failed to expand path (%s), error: %s", podfileParser.podfilePth, err)
	}

	envs := []string{fmt.Sprintf("PODFILE_PATH=%s", absPodfilePth)}
	podfileDir := filepath.Dir(absPodfilePth)

	out, err := runRubyScriptForOutput(rubyScriptContent, gemfileContent, podfileDir, envs)
	if err != nil {
		return map[string]string{}, fmt.Errorf("ruby script failed, error: %s", err)
	}

	if out == "" {
		return map[string]string{}, nil
	}

	type targetDefinitionOutputModel struct {
		Data  map[string]string
		Error string
	}

	var targetDefinitionOutput targetDefinitionOutputModel
	if err := json.Unmarshal([]byte(out), &targetDefinitionOutput); err != nil {
		return map[string]string{}, fmt.Errorf("failed to parse target definition output, error: %s", err)
	}

	if podfileParser.shouldRaiseReadDefinitionError(targetDefinitionOutput.Error) {
		return map[string]string{}, fmt.Errorf("failed to read target defintion map, error: %s", targetDefinitionOutput.Error)
	}

	return targetDefinitionOutput.Data, nil
}

func (podfileParser podfileParser) shouldRaiseReadDefinitionError(err string) bool {
	if err == "" {
		return false
	}

	isInvalidPodfileError := strings.Contains(err, "Pod::DSLError")
	if isInvalidPodfileError && podfileParser.suppressPodFileParseError {
		log.TWarnf("Could not parse podfile: %s", err)
		log.TWarnf("Will continue using default Cocoapods paths.")
		return false
	}

	return true
}

func (podfileParser podfileParser) getUserDefinedProjectRelavtivePath(cocoapodsVersion string) (string, error) {
	targetProjectMap, err := podfileParser.getTargetDefinitionProjectMap(cocoapodsVersion)
	if err != nil {
		return "", fmt.Errorf("failed to get target definition map, error: %s", err)
	}

	for target, project := range targetProjectMap {
		if target == "Pods" {
			return project, nil
		}
	}
	return "", nil
}

func (podfileParser podfileParser) getUserDefinedWorkspaceRelativePath(cocoapodsVersion string) (string, error) {
	gemfileCocoapodsVersion := ""
	if cocoapodsVersion != "" {
		gemfileCocoapodsVersion = fmt.Sprintf(`, '%s'`, cocoapodsVersion)
	}

	gemfileContent := fmt.Sprintf(`source 'https://rubygems.org'
gem 'cocoapods'%s
gem 'cocoapods-core'
gem 'json'
`, gemfileCocoapodsVersion)

	// returns WORKSPACE_NAME.xcworkspace if user defined a workspace name
	// returns empty struct {}, if no user defined workspace name exists in Podfile
	rubyScriptContent := `require 'cocoapods-core'
require 'json'

begin
	podfile_path = ENV['PODFILE_PATH']
	# In case of relative require in the Podfile, change working directory
	# For example: require_relative '../node_modules/react-native/scripts/react_native_pods'
	Dir.chdir(File.dirname(podfile_path))
	podfile = Pod::Podfile.from_file(podfile_path)
	pth = podfile.workspace_path
	puts "#{{ :data => pth }.to_json}"
rescue => e
	puts "#{{ :error => "#{e.to_s} Reason: #{e.message}"}.to_json}"
end
`
	absPodfilePth, err := filepath.Abs(podfileParser.podfilePth)
	if err != nil {
		return "", fmt.Errorf("failed to expand path (%s), error: %s", podfileParser.podfilePth, err)
	}

	envs := []string{fmt.Sprintf("PODFILE_PATH=%s", absPodfilePth)}
	podfileDir := filepath.Dir(absPodfilePth)

	out, err := runRubyScriptForOutput(rubyScriptContent, gemfileContent, podfileDir, envs)
	if err != nil {
		return "", fmt.Errorf("ruby script failed, error: %s", err)
	}

	if out == "" {
		return "", nil
	}

	type workspacePathOutputModel struct {
		Data  string
		Error string
	}

	var workspacePathOutput workspacePathOutputModel
	if err := json.Unmarshal([]byte(out), &workspacePathOutput); err != nil {
		return "", fmt.Errorf("failed to parse workspace path output, error: %s", err)
	}

	if podfileParser.shouldRaiseReadDefinitionError(workspacePathOutput.Error) {
		return "", fmt.Errorf("failed to read workspace path, error: %s", workspacePathOutput.Error)
	}

	return workspacePathOutput.Data, nil
}

// GetWorkspaceProjectMap ...
// If one project exists in the Podfile's directory, workspace name will be the project's name.
// If more then one project exists in the Podfile's directory, root 'xcodeproj/project' property have to be defined in the Podfile.
// Root 'xcodeproj/project' property will be mapped to the default cocoapods target (Pods).
// If workspace property defined in the Podfile, it will override the workspace name.
func (podfileParser podfileParser) GetWorkspaceProjectMap(projects []string) (map[string]string, error) {
	podfileDir := filepath.Dir(podfileParser.podfilePth)

	podfileLockPth, err := podfileParser.podfilelockPath(podfileDir)
	if err != nil {
		return map[string]string{}, err
	}

	cocoapodsVersion, err := podfileParser.cocoapodsVersion(podfileLockPth)
	if err != nil {
		return map[string]string{}, err
	}

	if err := podfileParser.fixPodfileQuotation(podfileParser.podfilePth); err != nil {
		return map[string]string{}, err
	}

	projectRelPth, err := podfileParser.getUserDefinedProjectRelavtivePath(cocoapodsVersion)
	if err != nil {
		return map[string]string{}, fmt.Errorf("failed to get user defined project path, error: %s", err)
	}

	if projectRelPth == "" {
		projects, err := pathutil.FilterPaths(projects, pathutil.InDirectoryFilter(podfileDir, true))
		if err != nil {
			return map[string]string{}, fmt.Errorf("failed to filter projects, error: %s", err)
		}

		if len(projects) == 0 {
			return map[string]string{}, errors.New("failed to determin workspace - project mapping: no explicit project specified and no project found in the Podfile's directory")
		} else if len(projects) > 1 {
			return map[string]string{}, errors.New("failed to determin workspace - project mapping: no explicit project specified and more than one project found in the Podfile's directory")
		}

		projectRelPth = filepath.Base(projects[0])
	}
	projectPth := filepath.Join(podfileDir, projectRelPth)

	if exist, err := pathutil.IsPathExists(projectPth); err != nil {
		return map[string]string{}, fmt.Errorf("failed to check if path (%s) exists, error: %s", projectPth, err)
	} else if !exist {
		return map[string]string{}, fmt.Errorf("project not found at: %s", projectPth)
	}

	workspaceRelPth, err := podfileParser.getUserDefinedWorkspaceRelativePath(cocoapodsVersion)
	if err != nil {
		return map[string]string{}, fmt.Errorf("failed to get user defined workspace path, error: %s", err)
	}

	if workspaceRelPth == "" {
		projectName := filepath.Base(strings.TrimSuffix(projectPth, ".xcodeproj"))
		workspaceRelPth = projectName + ".xcworkspace"
	}
	workspacePth := filepath.Join(podfileDir, workspaceRelPth)

	return map[string]string{
		workspacePth: projectPth,
	}, nil
}

func (podfileParser podfileParser) podfilelockPath(podfileDir string) (string, error) {
	podfileLockPth := filepath.Join(podfileDir, "Podfile.lock")
	if exist, err := pathutil.IsPathExists(podfileLockPth); err != nil {
		return "", fmt.Errorf("failed to check if Podfile.lock exist: %s", err)
	} else if !exist {
		podfileLockPth = filepath.Join(podfileDir, "podfile.lock")
		if exist, err := pathutil.IsPathExists(podfileLockPth); err != nil {
			return "", fmt.Errorf("failed to check if podfile.lock exist: %s", err)
		} else if !exist {
			podfileLockPth = ""
		}
	}

	return podfileLockPth, nil
}

func (podfileParser podfileParser) cocoapodsVersion(podfileLockPth string) (string, error) {

	if podfileLockPth == "" {
		return "", nil
	}

	version, err := GemVersionFromGemfileLock("cocoapods", podfileLockPth)
	if err != nil {
		return "", fmt.Errorf("failed to read cocoapods version from %s: %s", podfileLockPth, err)
	}
	return version, nil
}

func (podfileParser podfileParser) fixPodfileQuotation(podfilePth string) error {
	podfileContent, err := fileutil.ReadStringFromFile(podfilePth)
	if err != nil {
		return fmt.Errorf("failed to read podfile (%s), error: %s", podfilePth, err)
	}

	podfileContent = strings.Replace(podfileContent, `‘`, `'`, -1)
	podfileContent = strings.Replace(podfileContent, `’`, `'`, -1)
	podfileContent = strings.Replace(podfileContent, `“`, `"`, -1)
	podfileContent = strings.Replace(podfileContent, `”`, `"`, -1)

	if err := fileutil.WriteStringToFile(podfilePth, podfileContent); err != nil {
		return fmt.Errorf("failed to apply Podfile quotation fix, error: %s", err)
	}

	return nil
}

// MergePodWorkspaceProjectMap ...
// Previously we separated standalone projects and workspaces.
// But pod workspace-project map may define workspace which is not in the repository, but will be created by `pod install`.
// Related project should be found in the standalone projects list.
// We will create this workspace model, join the related project and remove this project from standlone projects.
// If workspace is in the repository, both workspace and project should be find in the input lists.
func MergePodWorkspaceProjectMap(podWorkspaceProjectMap map[string]string, standaloneProjects []xcodeproj.ProjectModel, workspaces []xcodeproj.WorkspaceModel) ([]xcodeproj.ProjectModel, []xcodeproj.WorkspaceModel, error) {
	mergedStandaloneProjects := []xcodeproj.ProjectModel{}
	mergedWorkspaces := []xcodeproj.WorkspaceModel{}

	for podWorkspaceFile, podProjectFile := range podWorkspaceProjectMap {
		podWorkspace, found := FindWorkspaceInList(podWorkspaceFile, workspaces)
		if found {
			// Workspace found, this means workspace is in the repository.
			podWorkspace.IsPodWorkspace = true

			// This case the project is already attached to the workspace.
			_, found := FindProjectInList(podProjectFile, podWorkspace.Projects)
			if !found {
				return []xcodeproj.ProjectModel{}, []xcodeproj.WorkspaceModel{}, fmt.Errorf("pod workspace (%s) found, but assigned project (%s) project not", podWorkspaceFile, podProjectFile)
			}

			// And the project is not standalone.
			_, found = FindProjectInList(podProjectFile, standaloneProjects)
			if found {
				return []xcodeproj.ProjectModel{}, []xcodeproj.WorkspaceModel{}, fmt.Errorf("pod workspace (%s) found, but assigned project (%s) marked as standalone", podWorkspaceFile, podProjectFile)
			}

			mergedStandaloneProjects = standaloneProjects
			mergedWorkspaces = ReplaceWorkspaceInList(workspaces, podWorkspace)
		} else {
			// Workspace not found, this means workspace is not in the repository,
			// but it will created by `pod install`.
			podWorkspace = xcodeproj.WorkspaceModel{
				Pth:            podWorkspaceFile,
				Name:           strings.TrimSuffix(filepath.Base(podWorkspaceFile), filepath.Ext(podWorkspaceFile)),
				IsPodWorkspace: true,
			}

			// This case the pod project was marked previously as standalone project.
			podProject, found := FindProjectInList(podProjectFile, standaloneProjects)
			if !found {
				return []xcodeproj.ProjectModel{}, []xcodeproj.WorkspaceModel{}, fmt.Errorf("pod workspace (%s) will be generated by (%s) project, but it does not found", podWorkspaceFile, podProjectFile)
			}

			podWorkspace.Projects = []xcodeproj.ProjectModel{podProject}

			mergedStandaloneProjects = RemoveProjectFromList(podProjectFile, standaloneProjects)
			mergedWorkspaces = append(workspaces, podWorkspace)
		}
	}

	return mergedStandaloneProjects, mergedWorkspaces, nil
}
