package ios

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-core/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xcode/xcodeproj"
	"github.com/stretchr/testify/require"
)

func TestAllowPodfileBaseFilter(t *testing.T) {
	t.Log("abs path")
	{
		absPaths := []string{
			"/Users/bitrise/Test.txt",
			"/Users/bitrise/.git/Podfile",
			"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Podfile",
			"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Podfile",
			"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Podfile",
		}

		expectedPaths := []string{
			"/Users/bitrise/.git/Podfile",
			"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Podfile",
			"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Podfile",
			"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Podfile",
		}

		actualPaths, err := utility.FilterPaths(absPaths, AllowPodfileBaseFilter)
		require.NoError(t, err)
		require.Equal(t, expectedPaths, actualPaths)
	}

	t.Log("rel path")
	{
		relPaths := []string{
			".",
			"Test.txt",
			".git/Podfile",
			"sample-apps-ios-cocoapods/Pods/Podfile",
			"ios-no-shared-schemes/Carthage/Checkouts/Result/Podfile",
			"ios-no-shared-schemes/test.framework/Checkouts/Result/Podfile",
		}

		expectedPaths := []string{
			".git/Podfile",
			"sample-apps-ios-cocoapods/Pods/Podfile",
			"ios-no-shared-schemes/Carthage/Checkouts/Result/Podfile",
			"ios-no-shared-schemes/test.framework/Checkouts/Result/Podfile",
		}

		actualPaths, err := utility.FilterPaths(relPaths, AllowPodfileBaseFilter)
		require.NoError(t, err)
		require.Equal(t, expectedPaths, actualPaths)
	}
}

func TestGetTargetDefinitionProjectMap(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__utility_test__")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	}()

	t.Log("xcodeproj defined")
	{
		tmpDir = filepath.Join(tmpDir, "xcodeproj_defined")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
project 'MyXcodeProject'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		expectedTargetDefinition := map[string]string{
			"Pods": "MyXcodeProject.xcodeproj",
		}
		actualTargetDefinition, err := getTargetDefinitionProjectMap(podfilePth, "")
		require.NoError(t, err)
		require.Equal(t, expectedTargetDefinition, actualTargetDefinition)
	}

	t.Log("xcodeproj NOT defined")
	{
		tmpDir = filepath.Join(tmpDir, "xcodeproj_not_defined")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		expectedTargetDefinition := map[string]string{}
		actualTargetDefinition, err := getTargetDefinitionProjectMap(podfilePth, "")
		require.NoError(t, err)
		require.Equal(t, expectedTargetDefinition, actualTargetDefinition)
	}

	t.Log("cocoapods 0.38.0")
	{
		tmpDir = filepath.Join(tmpDir, "xcodeproj_not_defined")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `source 'https://github.com/CocoaPods/Specs.git'
platform :ios, '8.0'

# pod 'Functional.m', '~> 1.0'

# Add Kiwi as an exclusive dependency for the Test target
target :SampleAppWithCocoapodsTests, :exclusive => true do
  pod 'Kiwi'
end

# post_install do |installer_representation|
#   installer_representation.project.targets.each do |target|
#     target.build_configurations.each do |config|
#       config.build_settings['ONLY_ACTIVE_ARCH'] = 'NO'
#     end
#   end
# end`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		expectedTargetDefinition := map[string]string{}
		actualTargetDefinition, err := getTargetDefinitionProjectMap(podfilePth, "0.38.0")
		require.NoError(t, err)
		require.Equal(t, expectedTargetDefinition, actualTargetDefinition)
	}
}

func TestGetUserDefinedProjectRelavtivePath(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__utility_test__")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	}()

	t.Log("xcodeproj defined")
	{
		tmpDir = filepath.Join(tmpDir, "xcodeproj_defined")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
project 'MyXcodeProject'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		expectedProject := "MyXcodeProject.xcodeproj"
		actualProject, err := getUserDefinedProjectRelavtivePath(podfilePth, "")
		require.NoError(t, err)
		require.Equal(t, expectedProject, actualProject)
	}

	t.Log("xcodeproj NOT defined")
	{
		tmpDir = filepath.Join(tmpDir, "xcodeproj_not_defined")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		expectedProject := ""
		actualProject, err := getUserDefinedProjectRelavtivePath(podfilePth, "")
		require.NoError(t, err)
		require.Equal(t, expectedProject, actualProject)
	}
}

func TestGetUserDefinedWorkspaceRelativePath(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__utility_test__")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	}()

	t.Log("workspace defined")
	{
		tmpDir = filepath.Join(tmpDir, "workspace_defined")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
workspace 'MyWorkspace'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		expectedWorkspace := "MyWorkspace.xcworkspace"
		actualWorkspace, err := getUserDefinedWorkspaceRelativePath(podfilePth, "")
		require.NoError(t, err)
		require.Equal(t, expectedWorkspace, actualWorkspace)
	}

	t.Log("workspace NOT defined")
	{
		tmpDir = filepath.Join(tmpDir, "workspace_not_defined")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		expectedWorkspace := ""
		actualWorkspace, err := getUserDefinedWorkspaceRelativePath(podfilePth, "")
		require.NoError(t, err)
		require.Equal(t, expectedWorkspace, actualWorkspace)
	}
}

func TestGetWorkspaceProjectMap(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__utility_test__")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	}()

	t.Log("0 project in Podfile's dir")
	{
		tmpDir = filepath.Join(tmpDir, "no_project")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		workspaceProjectMap, err := GetWorkspaceProjectMap(podfilePth, []string{})
		require.Error(t, err)
		require.Equal(t, 0, len(workspaceProjectMap))

		require.NoError(t, os.RemoveAll(tmpDir))
	}

	t.Log("1 project in Podfile's dir")
	{
		tmpDir = filepath.Join(tmpDir, "one_project")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		project := ""
		projectPth := filepath.Join(tmpDir, "project.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(projectPth, project))

		workspaceProjectMap, err := GetWorkspaceProjectMap(podfilePth, []string{projectPth})
		require.NoError(t, err)
		require.Equal(t, 1, len(workspaceProjectMap))

		for workspace, project := range workspaceProjectMap {
			workspaceBasename := filepath.Base(workspace)
			workspaceName := strings.TrimSuffix(workspaceBasename, ".xcworkspace")

			projectBasename := filepath.Base(project)
			projectName := strings.TrimSuffix(projectBasename, ".xcodeproj")

			require.Equal(t, "project", workspaceName, fmt.Sprintf("%v", workspaceProjectMap))
			require.Equal(t, "project", projectName, fmt.Sprintf("%v", workspaceProjectMap))
		}

		require.NoError(t, os.RemoveAll(tmpDir))
	}

	t.Log("Multiple project in Podfile's dir")
	{
		tmpDir = filepath.Join(tmpDir, "multiple_project")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		project1 := ""
		project1Pth := filepath.Join(tmpDir, "project1.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(project1Pth, project1))

		project2 := ""
		project2Pth := filepath.Join(tmpDir, "project2.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(project2Pth, project2))

		workspaceProjectMap, err := GetWorkspaceProjectMap(podfilePth, []string{project1Pth, project2Pth})
		require.Error(t, err)
		require.Equal(t, 0, len(workspaceProjectMap))

		require.NoError(t, os.RemoveAll(tmpDir))
	}

	t.Log("0 project in Podfile's dir + project defined in Podfile")
	{
		tmpDir = filepath.Join(tmpDir, "no_project_project_defined")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
project 'MyXcodeProject'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		workspaceProjectMap, err := GetWorkspaceProjectMap(podfilePth, []string{})
		require.Error(t, err)
		require.Equal(t, 0, len(workspaceProjectMap))

		require.NoError(t, os.RemoveAll(tmpDir))
	}

	t.Log("1 project in Podfile's dir + project defined in Podfile")
	{
		tmpDir = filepath.Join(tmpDir, "one_project_project_defined")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
project 'project'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		project := ""
		projectPth := filepath.Join(tmpDir, "project.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(projectPth, project))

		workspaceProjectMap, err := GetWorkspaceProjectMap(podfilePth, []string{projectPth})
		require.NoError(t, err)
		require.Equal(t, 1, len(workspaceProjectMap))

		for workspace, project := range workspaceProjectMap {
			workspaceBasename := filepath.Base(workspace)
			workspaceName := strings.TrimSuffix(workspaceBasename, ".xcworkspace")

			projectBasename := filepath.Base(project)
			projectName := strings.TrimSuffix(projectBasename, ".xcodeproj")

			require.Equal(t, "project", workspaceName, fmt.Sprintf("%v", workspaceProjectMap))
			require.Equal(t, "project", projectName, fmt.Sprintf("%v", workspaceProjectMap))
		}

		require.NoError(t, os.RemoveAll(tmpDir))
	}

	t.Log("Multiple project in Podfile's dir + project defined in Podfile")
	{
		tmpDir = filepath.Join(tmpDir, "multiple_project")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
project 'project1'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		project1 := ""
		project1Pth := filepath.Join(tmpDir, "project1.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(project1Pth, project1))

		project2 := ""
		project2Pth := filepath.Join(tmpDir, "project2.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(project2Pth, project2))

		workspaceProjectMap, err := GetWorkspaceProjectMap(podfilePth, []string{project1Pth, project2Pth})
		require.NoError(t, err)
		require.Equal(t, 1, len(workspaceProjectMap))

		for workspace, project := range workspaceProjectMap {
			workspaceBasename := filepath.Base(workspace)
			workspaceName := strings.TrimSuffix(workspaceBasename, ".xcworkspace")

			projectBasename := filepath.Base(project)
			projectName := strings.TrimSuffix(projectBasename, ".xcodeproj")

			require.Equal(t, "project1", workspaceName, fmt.Sprintf("%v", workspaceProjectMap))
			require.Equal(t, "project1", projectName, fmt.Sprintf("%v", workspaceProjectMap))
		}

		require.NoError(t, os.RemoveAll(tmpDir))
	}

	t.Log("1 project in Podfile's dir + workspace defined in Podfile")
	{
		tmpDir = filepath.Join(tmpDir, "one_project")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
workspace 'MyWorkspace'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		project := ""
		projectPth := filepath.Join(tmpDir, "project.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(projectPth, project))

		workspaceProjectMap, err := GetWorkspaceProjectMap(podfilePth, []string{projectPth})
		require.NoError(t, err)
		require.Equal(t, 1, len(workspaceProjectMap))

		for workspace, project := range workspaceProjectMap {
			workspaceBasename := filepath.Base(workspace)
			workspaceName := strings.TrimSuffix(workspaceBasename, ".xcworkspace")

			projectBasename := filepath.Base(project)
			projectName := strings.TrimSuffix(projectBasename, ".xcodeproj")

			require.Equal(t, "MyWorkspace", workspaceName, fmt.Sprintf("%v", workspaceProjectMap))
			require.Equal(t, "project", projectName, fmt.Sprintf("%v", workspaceProjectMap))
		}

		require.NoError(t, os.RemoveAll(tmpDir))
	}

	t.Log("Multiple project in Podfile's dir + workspace defined in Podfile")
	{
		tmpDir = filepath.Join(tmpDir, "multiple_project")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
project 'project1'
workspace 'MyWorkspace'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		project1 := ""
		project1Pth := filepath.Join(tmpDir, "project1.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(project1Pth, project1))

		project2 := ""
		project2Pth := filepath.Join(tmpDir, "project2.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(project2Pth, project2))

		workspaceProjectMap, err := GetWorkspaceProjectMap(podfilePth, []string{project1Pth, project2Pth})
		require.NoError(t, err)
		require.Equal(t, 1, len(workspaceProjectMap))

		for workspace, project := range workspaceProjectMap {
			workspaceBasename := filepath.Base(workspace)
			workspaceName := strings.TrimSuffix(workspaceBasename, ".xcworkspace")

			projectBasename := filepath.Base(project)
			projectName := strings.TrimSuffix(projectBasename, ".xcodeproj")

			require.Equal(t, "MyWorkspace", workspaceName, fmt.Sprintf("%v", workspaceProjectMap))
			require.Equal(t, "project1", projectName, fmt.Sprintf("%v", workspaceProjectMap))
		}

		require.NoError(t, os.RemoveAll(tmpDir))
	}
}

func TestMergePodWorkspaceProjectMap(t *testing.T) {
	t.Log("workspace is in the repository")
	{
		podWorkspaceMap := map[string]string{
			"MyWorkspace.xcworkspace": "MyXcodeProject.xcodeproj",
		}

		standaloneProjects := []xcodeproj.ProjectModel{}
		expectedStandaloneProjects := []xcodeproj.ProjectModel{}

		workspaces := []xcodeproj.WorkspaceModel{
			xcodeproj.WorkspaceModel{
				Pth:  "MyWorkspace.xcworkspace",
				Name: "MyWorkspace",
				Projects: []xcodeproj.ProjectModel{
					xcodeproj.ProjectModel{
						Pth: "MyXcodeProject.xcodeproj",
					},
				},
			},
		}
		expectedWorkspaces := []xcodeproj.WorkspaceModel{
			xcodeproj.WorkspaceModel{
				Pth:  "MyWorkspace.xcworkspace",
				Name: "MyWorkspace",
				Projects: []xcodeproj.ProjectModel{
					xcodeproj.ProjectModel{
						Pth: "MyXcodeProject.xcodeproj",
					},
				},
				IsPodWorkspace: true,
			},
		}

		mergedStandaloneProjects, mergedWorkspaces, err := MergePodWorkspaceProjectMap(podWorkspaceMap, standaloneProjects, workspaces)
		require.NoError(t, err)
		require.Equal(t, expectedStandaloneProjects, mergedStandaloneProjects)
		require.Equal(t, expectedWorkspaces, mergedWorkspaces)
	}

	t.Log("workspace is in the repository, but project not attached - ERROR")
	{
		podWorkspaceMap := map[string]string{
			"MyWorkspace.xcworkspace": "MyXcodeProject.xcodeproj",
		}

		standaloneProjects := []xcodeproj.ProjectModel{}

		workspaces := []xcodeproj.WorkspaceModel{
			xcodeproj.WorkspaceModel{
				Pth:  "MyWorkspace.xcworkspace",
				Name: "MyWorkspace",
			},
		}

		mergedStandaloneProjects, mergedWorkspaces, err := MergePodWorkspaceProjectMap(podWorkspaceMap, standaloneProjects, workspaces)
		require.Error(t, err)
		require.Equal(t, []xcodeproj.ProjectModel{}, mergedStandaloneProjects)
		require.Equal(t, []xcodeproj.WorkspaceModel{}, mergedWorkspaces)
	}

	t.Log("workspace is in the repository, but project is marged as standalon - ERROR")
	{
		podWorkspaceMap := map[string]string{
			"MyWorkspace.xcworkspace": "MyXcodeProject.xcodeproj",
		}

		standaloneProjects := []xcodeproj.ProjectModel{
			xcodeproj.ProjectModel{
				Pth: "MyXcodeProject.xcodeproj",
			},
		}

		workspaces := []xcodeproj.WorkspaceModel{
			xcodeproj.WorkspaceModel{
				Pth:  "MyWorkspace.xcworkspace",
				Name: "MyWorkspace",
			},
		}

		mergedStandaloneProjects, mergedWorkspaces, err := MergePodWorkspaceProjectMap(podWorkspaceMap, standaloneProjects, workspaces)
		require.Error(t, err)
		require.Equal(t, []xcodeproj.ProjectModel{}, mergedStandaloneProjects)
		require.Equal(t, []xcodeproj.WorkspaceModel{}, mergedWorkspaces)
	}

	t.Log("workspace is gitignored")
	{
		podWorkspaceMap := map[string]string{
			"MyWorkspace.xcworkspace": "MyXcodeProject.xcodeproj",
		}

		standaloneProjects := []xcodeproj.ProjectModel{
			xcodeproj.ProjectModel{
				Pth: "MyXcodeProject.xcodeproj",
			},
		}
		expectedStandaloneProjects := []xcodeproj.ProjectModel{}

		workspaces := []xcodeproj.WorkspaceModel{}
		expectedWorkspaces := []xcodeproj.WorkspaceModel{
			xcodeproj.WorkspaceModel{
				Pth:  "MyWorkspace.xcworkspace",
				Name: "MyWorkspace",
				Projects: []xcodeproj.ProjectModel{
					xcodeproj.ProjectModel{
						Pth: "MyXcodeProject.xcodeproj",
					},
				},
				IsPodWorkspace: true,
			},
		}

		mergedStandaloneProjects, mergedWorkspaces, err := MergePodWorkspaceProjectMap(podWorkspaceMap, standaloneProjects, workspaces)
		require.NoError(t, err)
		require.Equal(t, expectedStandaloneProjects, mergedStandaloneProjects)
		require.Equal(t, expectedWorkspaces, mergedWorkspaces)
	}

	t.Log("workspace is gitignored, but standalon project missing - ERROR")
	{
		podWorkspaceMap := map[string]string{
			"MyWorkspace.xcworkspace": "MyXcodeProject.xcodeproj",
		}

		standaloneProjects := []xcodeproj.ProjectModel{}

		workspaces := []xcodeproj.WorkspaceModel{}

		mergedStandaloneProjects, mergedWorkspaces, err := MergePodWorkspaceProjectMap(podWorkspaceMap, standaloneProjects, workspaces)
		require.Error(t, err)
		require.Equal(t, []xcodeproj.ProjectModel{}, mergedStandaloneProjects)
		require.Equal(t, []xcodeproj.WorkspaceModel{}, mergedWorkspaces)
	}
}
