package builder

import (
	"fmt"

	"github.com/bitrise-tools/go-xamarin/analyzers/project"
	"github.com/bitrise-tools/go-xamarin/constants"
	"github.com/bitrise-tools/go-xamarin/tools"
	"github.com/bitrise-tools/go-xamarin/tools/buildtools"
	"github.com/bitrise-tools/go-xamarin/tools/buildtools/msbuild"
	"github.com/bitrise-tools/go-xamarin/tools/buildtools/xbuild"
	"github.com/bitrise-tools/go-xamarin/tools/nunit"
	"github.com/bitrise-tools/go-xamarin/utility"
)

func (builder Model) buildSolutionCommand(configuration, platform string) (tools.Runnable, error) {
	var buildCommand tools.Runnable

	var command *xbuild.Model
	var err error

	if builder.buildTool == buildtools.Msbuild {
		command, err = msbuild.New(builder.solution.Pth, "")
	} else {
		command, err = xbuild.New(builder.solution.Pth, "")
	}

	if err != nil {
		return nil, err
	}

	command.SetTarget("Build")
	command.SetConfiguration(configuration)
	command.SetPlatform(platform)
	buildCommand = command

	return buildCommand, nil
}

func (builder Model) buildProjectCommand(configuration, platform string, proj project.Model) ([]tools.Runnable, []string, error) {
	warnings := []string{}

	solutionConfig := utility.ToConfig(configuration, platform)

	projectConfigKey, ok := proj.ConfigMap[solutionConfig]
	if !ok {
		warnings = append(warnings, fmt.Sprintf("project (%s) do not have config for solution config (%s), skipping...", proj.Name, solutionConfig))
	}

	projectConfig, ok := proj.Configs[projectConfigKey]
	if !ok {
		warnings = append(warnings, fmt.Sprintf("project (%s) contains mapping for solution config (%s), but does not have project configuration", proj.Name, solutionConfig))
	}

	// Prepare build commands
	buildCommands := []tools.Runnable{}

	switch proj.SDK {
	case constants.SDKIOS, constants.SDKTvOS:
		var command *xbuild.Model
		var err error

		if builder.buildTool == buildtools.Msbuild {
			command, err = msbuild.New(builder.solution.Pth, "")
		} else {
			command, err = xbuild.New(builder.solution.Pth, "")
		}

		if err != nil {
			return []tools.Runnable{}, warnings, err
		}

		command.SetTarget("Build")
		command.SetConfiguration(configuration)
		command.SetPlatform(platform)

		if isArchitectureArchiveable(projectConfig.MtouchArchs...) {
			command.SetBuildIpa(true)
			command.SetArchiveOnBuild(true)
		}

		buildCommands = append(buildCommands, command)
	case constants.SDKMacOS:
		var command *xbuild.Model
		var err error

		if builder.buildTool == buildtools.Msbuild {
			command, err = msbuild.New(builder.solution.Pth, "")
		} else {
			command, err = xbuild.New(builder.solution.Pth, "")
		}

		if err != nil {
			return []tools.Runnable{}, warnings, err
		}

		command.SetTarget("Build")
		command.SetConfiguration(configuration)
		command.SetPlatform(platform)
		command.SetArchiveOnBuild(true)

		buildCommands = append(buildCommands, command)
	case constants.SDKAndroid:
		var command *xbuild.Model
		var err error

		if builder.buildTool == buildtools.Msbuild {
			command, err = msbuild.New(builder.solution.Pth, proj.Pth)
		} else {
			command, err = xbuild.New(builder.solution.Pth, proj.Pth)
		}

		if err != nil {
			return []tools.Runnable{}, warnings, err
		}

		if projectConfig.SignAndroid {
			command.SetTarget("SignAndroidPackage")
		} else {
			command.SetTarget("PackageForAndroid")
		}

		command.SetConfiguration(projectConfig.Configuration)

		if !isPlatformAnyCPU(projectConfig.Platform) {
			command.SetPlatform(projectConfig.Platform)
		}

		buildCommands = append(buildCommands, command)
	}

	return buildCommands, warnings, nil
}

func (builder Model) buildXamarinUITestProjectCommand(configuration, platform string, proj project.Model) (tools.Runnable, []string, error) {
	warnings := []string{}

	solutionConfig := utility.ToConfig(configuration, platform)

	projectConfigKey, ok := proj.ConfigMap[solutionConfig]
	if !ok {
		warnings = append(warnings, fmt.Sprintf("project (%s) do not have config for solution config (%s), skipping...", proj.Name, solutionConfig))
	}

	projectConfig, ok := proj.Configs[projectConfigKey]
	if !ok {
		warnings = append(warnings, fmt.Sprintf("project (%s) contains mapping for solution config (%s), but does not have project configuration", proj.Name, solutionConfig))
	}

	var command *xbuild.Model
	var err error

	if builder.buildTool == buildtools.Msbuild {
		command, err = msbuild.New(builder.solution.Pth, proj.Pth)
	} else {
		command, err = xbuild.New(builder.solution.Pth, proj.Pth)
	}
	if err != nil {
		return nil, warnings, err
	}

	command.SetConfiguration(projectConfig.Configuration)

	return command, warnings, nil
}

func (builder Model) buildNunitTestProjectCommand(configuration, platform string, proj project.Model, nunitConsolePth string) (tools.Runnable, []string, error) {
	warnings := []string{}

	solutionConfig := utility.ToConfig(configuration, platform)

	projectConfigKey, ok := proj.ConfigMap[solutionConfig]
	if !ok {
		warnings = append(warnings, fmt.Sprintf("project (%s) do not have config for solution config (%s), skipping...", proj.Name, solutionConfig))
	}

	projectConfig, ok := proj.Configs[projectConfigKey]
	if !ok {
		warnings = append(warnings, fmt.Sprintf("project (%s) contains mapping for solution config (%s), but does not have project configuration", proj.Name, solutionConfig))
	}

	command, err := nunit.New(nunitConsolePth)
	if err != nil {
		return nil, warnings, err
	}

	command.SetProjectPth(proj.Pth)
	command.SetConfig(projectConfig.Configuration)

	return command, warnings, nil
}
