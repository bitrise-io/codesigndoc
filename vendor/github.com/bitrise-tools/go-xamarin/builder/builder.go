package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xamarin/analyzers/project"
	"github.com/bitrise-tools/go-xamarin/analyzers/solution"
	"github.com/bitrise-tools/go-xamarin/constants"
	"github.com/bitrise-tools/go-xamarin/tools"
	"github.com/bitrise-tools/go-xamarin/tools/buildtools"
	"github.com/bitrise-tools/go-xamarin/tools/nunit"
	"github.com/bitrise-tools/go-xamarin/utility"
)

// Model ...
type Model struct {
	solution solution.Model

	projectTypeWhitelist []constants.SDK
	buildTool            buildtools.BuildTool
}

// OutputModel ...
type OutputModel struct {
	Pth        string
	OutputType constants.OutputType
}

// ProjectOutputModel ...
type ProjectOutputModel struct {
	ProjectType constants.SDK
	Outputs     []OutputModel
}

// ProjectOutputMap ...
type ProjectOutputMap map[string]ProjectOutputModel // Project Name - ProjectOutputModel

// TestProjectOutputModel ...
type TestProjectOutputModel struct {
	TestFramwork         constants.TestFramework
	ReferredProjectNames []string
	Output               OutputModel
}

// TestProjectOutputMap ...
type TestProjectOutputMap map[string]TestProjectOutputModel // Test Project Name - TestProjectOutputModel

// PrepareCommandCallback ...
type PrepareCommandCallback func(solutionName string, projectName string, sdk constants.SDK, testFramework constants.TestFramework, command *tools.Editable)

// BuildCommandCallback ...
type BuildCommandCallback func(solutionName string, projectName string, sdk constants.SDK, testFramework constants.TestFramework, commandStr string, alreadyPerformed bool)

// ClearCommandCallback ...
type ClearCommandCallback func(project project.Model, dir string)

// New ...
func New(solutionPth string, projectTypeWhitelist []constants.SDK, buildTool buildtools.BuildTool) (Model, error) {
	if err := validateSolutionPth(solutionPth); err != nil {
		return Model{}, err
	}

	solution, err := solution.New(solutionPth, true)
	if err != nil {
		return Model{}, err
	}

	if projectTypeWhitelist == nil {
		projectTypeWhitelist = []constants.SDK{}
	}

	return Model{
		solution: solution,

		projectTypeWhitelist: projectTypeWhitelist,
		buildTool:            buildTool,
	}, nil
}

// CleanAll ...
func (builder Model) CleanAll(callback ClearCommandCallback) error {
	whitelistedProjects := builder.whitelistedProjects()

	for _, proj := range whitelistedProjects {

		projectDir := filepath.Dir(proj.Pth)

		{
			binPth := filepath.Join(projectDir, "bin")
			if exist, err := pathutil.IsDirExists(binPth); err != nil {
				return err
			} else if exist {
				if callback != nil {
					callback(proj, binPth)
				}

				if err := os.RemoveAll(binPth); err != nil {
					return err
				}
			}
		}

		{
			objPth := filepath.Join(projectDir, "obj")
			if exist, err := pathutil.IsDirExists(objPth); err != nil {
				return err
			} else if exist {
				if callback != nil {
					callback(proj, objPth)
				}

				if err := os.RemoveAll(objPth); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// BuildSolution ...
func (builder Model) BuildSolution(configuration, platform string, callback BuildCommandCallback) error {
	if err := validateSolutionConfig(builder.solution, configuration, platform); err != nil {
		return err
	}

	buildCommand, err := builder.buildSolutionCommand(configuration, platform)
	if err != nil {
		return fmt.Errorf("Failed to create build command, error: %s", err)
	}

	// Callback to notify the caller about next running command
	if callback != nil {
		callback(builder.solution.Name, "", constants.SDKUnknown, constants.TestFrameworkUnknown, buildCommand.PrintableCommand(), false)
	}

	return buildCommand.Run()
}

// BuildAllProjects ...
func (builder Model) BuildAllProjects(configuration, platform string, prepareCallback PrepareCommandCallback, callback BuildCommandCallback) ([]string, error) {
	warnings := []string{}

	if err := validateSolutionConfig(builder.solution, configuration, platform); err != nil {
		return warnings, err
	}

	buildableProjects, warns := builder.buildableProjects(configuration, platform)
	if len(buildableProjects) == 0 {
		return warns, fmt.Errorf("No project to build found")
	}

	perfomedCommands := []tools.Printable{}

	for _, proj := range buildableProjects {
		buildCommands, warns, err := builder.buildProjectCommand(configuration, platform, proj)
		warnings = append(warnings, warns...)
		if err != nil {
			return warnings, fmt.Errorf("Failed to create build command, error: %s", err)
		}

		for _, buildCommand := range buildCommands {
			// Callback to let the caller to modify the command
			if prepareCallback != nil {
				editabeCommand := tools.Editable(buildCommand)
				prepareCallback(builder.solution.Name, proj.Name, proj.SDK, proj.TestFramework, &editabeCommand)
			}

			// Check if same command was already performed
			alreadyPerformed := false
			if tools.PrintableSliceContains(perfomedCommands, buildCommand) {
				alreadyPerformed = true
			}

			// Callback to notify the caller about next running command
			if callback != nil {
				callback(builder.solution.Name, proj.Name, proj.SDK, proj.TestFramework, buildCommand.PrintableCommand(), alreadyPerformed)
			}

			if !alreadyPerformed {
				if err := buildCommand.Run(); err != nil {
					return warnings, err
				}
				perfomedCommands = append(perfomedCommands, buildCommand)
			}
		}
	}

	return warnings, nil
}

// BuildAllUITestableXamarinProjects ...
func (builder Model) BuildAllUITestableXamarinProjects(configuration, platform string, prepareCallback PrepareCommandCallback, callback BuildCommandCallback) ([]string, error) {
	warnings := []string{}

	if err := validateSolutionConfig(builder.solution, configuration, platform); err != nil {
		return warnings, err
	}

	if err := builder.BuildSolution(configuration, platform, callback); err != nil {
		return nil, err
	}

	_, buildableReferredProjects, warns := builder.buildableXamarinUITestProjectsAndReferredProjects(configuration, platform)
	if len(buildableReferredProjects) == 0 {
		return warns, fmt.Errorf("No project to build found")
	}

	perfomedCommands := []tools.Printable{}

	for _, proj := range buildableReferredProjects {
		buildCommands, warns, err := builder.buildProjectCommand(configuration, platform, proj)
		warnings = append(warnings, warns...)
		if err != nil {
			return warnings, fmt.Errorf("Failed to create build command, error: %s", err)
		}

		for _, buildCommand := range buildCommands {
			// Callback to let the caller to modify the command
			if prepareCallback != nil {
				editabeCommand := tools.Editable(buildCommand)
				prepareCallback(builder.solution.Name, proj.Name, proj.SDK, proj.TestFramework, &editabeCommand)
			}

			// Check if same command was already performed
			alreadyPerformed := false
			if tools.PrintableSliceContains(perfomedCommands, buildCommand) {
				alreadyPerformed = true
			}

			// Callback to notify the caller about next running command
			if callback != nil {
				callback(builder.solution.Name, proj.Name, proj.SDK, proj.TestFramework, buildCommand.PrintableCommand(), alreadyPerformed)
			}

			if !alreadyPerformed {
				if err := buildCommand.Run(); err != nil {
					return warnings, err
				}
				perfomedCommands = append(perfomedCommands, buildCommand)
			}
		}
	}

	return warnings, nil
}

// RunAllXamarinUITests ...
func (builder Model) RunAllXamarinUITests(configuration, platform string, prepareCallback PrepareCommandCallback, callback BuildCommandCallback) ([]string, error) {
	warnings := []string{}

	if err := validateSolutionConfig(builder.solution, configuration, platform); err != nil {
		return warnings, err
	}

	buildableTestProjects, _, warns := builder.buildableXamarinUITestProjectsAndReferredProjects(configuration, platform)
	if len(buildableTestProjects) == 0 {
		return warns, fmt.Errorf("No project to build found")
	}

	perfomedCommands := []tools.Printable{}

	for _, testProj := range buildableTestProjects {
		buildCommand, warns, err := builder.buildXamarinUITestProjectCommand(configuration, platform, testProj)
		warnings = append(warnings, warns...)
		if err != nil {
			return warnings, fmt.Errorf("Failed to create build command, error: %s", err)
		}

		// Callback to let the caller to modify the command
		if prepareCallback != nil {
			editabeCommand := tools.Editable(buildCommand)
			prepareCallback(builder.solution.Name, testProj.Name, testProj.SDK, testProj.TestFramework, &editabeCommand)
		}

		// Check if same command was already performed
		alreadyPerformed := false
		if tools.PrintableSliceContains(perfomedCommands, buildCommand) {
			alreadyPerformed = true
		}

		// Callback to notify the caller about next running command
		if callback != nil {
			callback(builder.solution.Name, testProj.Name, testProj.SDK, testProj.TestFramework, buildCommand.PrintableCommand(), alreadyPerformed)
		}

		if !alreadyPerformed {
			if err := buildCommand.Run(); err != nil {
				return warnings, err
			}
			perfomedCommands = append(perfomedCommands, buildCommand)
		}
	}

	return warnings, nil
}

// BuildAndRunAllXamarinUITestAndReferredProjects ...
func (builder Model) BuildAndRunAllXamarinUITestAndReferredProjects(configuration, platform string, prepareCallback PrepareCommandCallback, callback BuildCommandCallback) ([]string, error) {
	warnings := []string{}

	buildWarnings, err := builder.BuildAllUITestableXamarinProjects(configuration, platform, prepareCallback, callback)
	warnings = append(warnings, buildWarnings...)
	if err != nil {
		return warnings, err
	}

	runWarnings, err := builder.RunAllXamarinUITests(configuration, platform, prepareCallback, callback)
	warnings = append(warnings, runWarnings...)
	if err != nil {
		return warnings, err
	}

	return warnings, nil
}

// RunAllNunitTestProjects ...
func (builder Model) RunAllNunitTestProjects(configuration, platform string, callback BuildCommandCallback, prepareCallback PrepareCommandCallback) ([]string, error) {
	if err := validateSolutionConfig(builder.solution, configuration, platform); err != nil {
		return nil, err
	}

	buildableProjects, warns := builder.buildableNunitTestProjects(configuration, platform)
	if len(buildableProjects) == 0 {
		return warns, fmt.Errorf("No project to build found")
	}

	nunitConsolePth, err := nunit.SystemNunit3ConsolePath()
	if err != nil {
		return nil, err
	}

	warnings := []string{}
	perfomedCommands := []tools.Printable{}

	for _, testProj := range buildableProjects {
		buildCommand, warns, err := builder.buildNunitTestProjectCommand(configuration, platform, testProj, nunitConsolePth)
		warnings = append(warnings, warns...)
		if err != nil {
			return warnings, fmt.Errorf("Failed to create build command, error: %s", err)
		}

		// Callback to let the caller to modify the command
		if prepareCallback != nil {
			editabeCommand := tools.Editable(buildCommand)
			prepareCallback(builder.solution.Name, testProj.Name, constants.SDKUnknown, constants.TestFrameworkNunitTest, &editabeCommand)
		}

		// Check if same command was already performed
		alreadyPerformed := false
		if tools.PrintableSliceContains(perfomedCommands, buildCommand) {
			alreadyPerformed = true
		}

		// Callback to notify the caller about next running command
		if callback != nil {
			callback(builder.solution.Name, testProj.Name, constants.SDKUnknown, constants.TestFrameworkNunitTest, buildCommand.PrintableCommand(), alreadyPerformed)
		}

		if !alreadyPerformed {
			if err := buildCommand.Run(); err != nil {
				return warnings, err
			}
			perfomedCommands = append(perfomedCommands, buildCommand)
		}
	}

	return warnings, nil
}

// BuildAndRunAllNunitTestProjects ...
func (builder Model) BuildAndRunAllNunitTestProjects(configuration, platform string, callback BuildCommandCallback, prepareCallback PrepareCommandCallback) ([]string, error) {
	if err := builder.BuildSolution(configuration, platform, callback); err != nil {
		return nil, err
	}

	return builder.RunAllNunitTestProjects(configuration, platform, callback, prepareCallback)
}

// CollectProjectOutputs ...
func (builder Model) CollectProjectOutputs(configuration, platform string, startTime, endTime time.Time) (ProjectOutputMap, error) {
	projectOutputMap := ProjectOutputMap{}

	buildableProjects, _ := builder.buildableProjects(configuration, platform)

	solutionConfig := utility.ToConfig(configuration, platform)

	for _, proj := range buildableProjects {
		projectConfigKey, ok := proj.ConfigMap[solutionConfig]
		if !ok {
			continue
		}

		projectConfig, ok := proj.Configs[projectConfigKey]
		if !ok {
			continue
		}

		projectOutputs, ok := projectOutputMap[proj.Name]
		if !ok {
			projectOutputs = ProjectOutputModel{
				ProjectType: proj.SDK,
				Outputs:     []OutputModel{},
			}
		}

		switch proj.SDK {
		case constants.SDKIOS, constants.SDKTvOS:
			if isArchitectureArchiveable(projectConfig.MtouchArchs...) {
				if xcarchivePth, err := exportLatestXCArchiveFromXcodeArchives(proj.AssemblyName, startTime, endTime); err != nil {
					return ProjectOutputMap{}, err
				} else if xcarchivePth != "" {
					projectOutputs.Outputs = append(projectOutputs.Outputs, OutputModel{
						Pth:        xcarchivePth,
						OutputType: constants.OutputTypeXCArchive,
					})
				}

				if ipaPth, err := exportLatestIpa(projectConfig.OutputDir, proj.AssemblyName, startTime, endTime); err != nil {
					return ProjectOutputMap{}, err
				} else if ipaPth != "" {
					projectOutputs.Outputs = append(projectOutputs.Outputs, OutputModel{
						Pth:        ipaPth,
						OutputType: constants.OutputTypeIPA,
					})
				}

				if dsymPth, err := exportAppDSYM(projectConfig.OutputDir, proj.AssemblyName, startTime, endTime); err != nil {
					return ProjectOutputMap{}, err
				} else if dsymPth != "" {
					projectOutputs.Outputs = append(projectOutputs.Outputs, OutputModel{
						Pth:        dsymPth,
						OutputType: constants.OutputTypeDSYM,
					})
				}
			}

			if appPth, err := exportApp(projectConfig.OutputDir, proj.AssemblyName, startTime, endTime); err != nil {
				return ProjectOutputMap{}, err
			} else if appPth != "" {
				projectOutputs.Outputs = append(projectOutputs.Outputs, OutputModel{
					Pth:        appPth,
					OutputType: constants.OutputTypeAPP,
				})
			}
		case constants.SDKMacOS:
			if appPth, err := exportApp(projectConfig.OutputDir, proj.AssemblyName, startTime, endTime); err != nil {
				return ProjectOutputMap{}, err
			} else if appPth != "" {
				projectOutputs.Outputs = append(projectOutputs.Outputs, OutputModel{
					Pth:        appPth,
					OutputType: constants.OutputTypeAPP,
				})
			}
			if pkgPth, err := exportPKG(projectConfig.OutputDir, proj.AssemblyName, startTime, endTime); err != nil {
				return ProjectOutputMap{}, err
			} else if pkgPth != "" {
				projectOutputs.Outputs = append(projectOutputs.Outputs, OutputModel{
					Pth:        pkgPth,
					OutputType: constants.OutputTypePKG,
				})
			}
		case constants.SDKAndroid:
			packageName, err := androidPackageName(proj.ManifestPth)
			if err != nil {
				return ProjectOutputMap{}, err
			}

			if apkPth, err := exportApk(projectConfig.OutputDir, packageName, startTime, endTime); err != nil {
				return ProjectOutputMap{}, err
			} else if apkPth != "" {
				projectOutputs.Outputs = append(projectOutputs.Outputs, OutputModel{
					Pth:        apkPth,
					OutputType: constants.OutputTypeAPK,
				})
			}
		}

		if len(projectOutputs.Outputs) > 0 {
			projectOutputMap[proj.Name] = projectOutputs
		}
	}

	return projectOutputMap, nil
}

// CollectXamarinUITestProjectOutputs ...
func (builder Model) CollectXamarinUITestProjectOutputs(configuration, platform string, startTime, endTime time.Time) (TestProjectOutputMap, []string, error) {
	testProjectOutputMap := TestProjectOutputMap{}
	warnings := []string{}

	buildableTestProjects, _, _ := builder.buildableXamarinUITestProjectsAndReferredProjects(configuration, platform)

	solutionConfig := utility.ToConfig(configuration, platform)

	for _, testProj := range buildableTestProjects {
		projectConfigKey, ok := testProj.ConfigMap[solutionConfig]
		if !ok {
			continue
		}

		projectConfig, ok := testProj.Configs[projectConfigKey]
		if !ok {
			continue
		}

		if dllPth, err := exportDLL(projectConfig.OutputDir, testProj.AssemblyName, startTime, endTime); err != nil {
			return TestProjectOutputMap{}, warnings, err
		} else if dllPth != "" {
			referredProjectNames := []string{}
			referredProjectIDs := testProj.ReferredProjectIDs
			for _, referredProjectID := range referredProjectIDs {
				referredProject, ok := builder.solution.ProjectMap[referredProjectID]
				if !ok {
					warnings = append(warnings, fmt.Sprintf("project reference exist with project id: %s, but project not found in solution", referredProjectID))
				}

				referredProjectNames = append(referredProjectNames, referredProject.Name)
			}

			testProjectOutputMap[testProj.Name] = TestProjectOutputModel{
				TestFramwork:         testProj.TestFramework,
				ReferredProjectNames: referredProjectNames,
				Output: OutputModel{
					Pth:        dllPth,
					OutputType: constants.OutputTypeDLL,
				},
			}
		}
	}

	return testProjectOutputMap, warnings, nil
}
