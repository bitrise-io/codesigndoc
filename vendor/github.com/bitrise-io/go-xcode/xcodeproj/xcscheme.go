package xcodeproj

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcscheme"
)

// SchemeModel ...
type SchemeModel struct {
	Name                  string
	HasXCTest             bool
	BuildableReferenceIDs []string
}

func filterSharedSchemeFilePaths(paths []string) []string {
	isSharedSchemeFilePath := func(pth string) bool {
		regexpPattern := filepath.Join(".*[/]?xcshareddata", "xcschemes", ".+[.]xcscheme")
		regexp := regexp.MustCompile(regexpPattern)
		return (regexp.FindString(pth) != "")
	}

	filteredPaths := []string{}
	for _, pth := range paths {
		if isSharedSchemeFilePath(pth) {
			filteredPaths = append(filteredPaths, pth)
		}
	}

	sort.Strings(filteredPaths)

	return filteredPaths
}

func sharedSchemeFilePaths(projectOrWorkspacePth string) ([]string, error) {
	filesInDir := func(dir string) ([]string, error) {
		files := []string{}
		if err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			files = append(files, path)
			return nil
		}); err != nil {
			return []string{}, err
		}
		return files, nil
	}

	paths, err := filesInDir(projectOrWorkspacePth)
	if err != nil {
		return []string{}, err
	}
	return filterSharedSchemeFilePaths(paths), nil
}

// SchemeNameFromPath ...
func SchemeNameFromPath(schemePth string) string {
	basename := filepath.Base(schemePth)
	ext := filepath.Ext(schemePth)
	if ext != XCSchemeExt {
		return ""
	}
	return strings.TrimSuffix(basename, ext)
}

func schemeFileContentContainsXCTestBuildAction(schemeFileContent string) (bool, error) {
	testActionStartPattern := "<TestAction"
	testActionEndPattern := "</TestAction>"
	isTestableAction := false

	testableReferenceStartPattern := "<TestableReference"
	testableReferenceSkippedRegexp := regexp.MustCompile(`skipped = "(?P<skipped>.+)"`)
	testableReferenceEndPattern := "</TestableReference>"
	isTestableReference := false

	xctestBuildableReferenceNameRegexp := regexp.MustCompile(`BuildableName = ".+.xctest"`)

	scanner := bufio.NewScanner(strings.NewReader(schemeFileContent))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == testActionEndPattern {
			break
		}

		if strings.TrimSpace(line) == testActionStartPattern {
			isTestableAction = true
			continue
		}

		if !isTestableAction {
			continue
		}

		// TestAction

		if strings.TrimSpace(line) == testableReferenceEndPattern {
			isTestableReference = false
			continue
		}

		if strings.TrimSpace(line) == testableReferenceStartPattern {
			isTestableReference = true
			continue
		}

		if !isTestableReference {
			continue
		}

		// TestableReference

		if matches := testableReferenceSkippedRegexp.FindStringSubmatch(line); len(matches) > 1 {
			skipped := matches[1]
			if skipped != "NO" {
				break
			}
		}

		if match := xctestBuildableReferenceNameRegexp.FindString(line); match != "" {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

// SchemeFileContainsXCTestBuildAction ...
func SchemeFileContainsXCTestBuildAction(schemeFilePth string) (bool, error) {
	content, err := fileutil.ReadStringFromFile(schemeFilePth)
	if err != nil {
		return false, err
	}

	return schemeFileContentContainsXCTestBuildAction(content)
}

func sharedSchemes(projectOrWorkspacePth string) ([]SchemeModel, error) {
	schemePaths, err := sharedSchemeFilePaths(projectOrWorkspacePth)
	if err != nil {
		return []SchemeModel{}, err
	}

	schemes := []SchemeModel{}
	for _, schemePth := range schemePaths {
		schemeName := SchemeNameFromPath(schemePth)

		hasXCTest, err := SchemeFileContainsXCTestBuildAction(schemePth)
		if err != nil {
			return []SchemeModel{}, err
		}

		buildableReferenceIDs, err := buildableReferenceIDs(schemePth)
		if err != nil {
			return []SchemeModel{}, err
		}

		schemes = append(schemes, SchemeModel{
			Name:                  schemeName,
			HasXCTest:             hasXCTest,
			BuildableReferenceIDs: buildableReferenceIDs,
		})
	}

	return schemes, nil
}

// ProjectSharedSchemes ...
func ProjectSharedSchemes(projectPth string) ([]SchemeModel, error) {
	return sharedSchemes(projectPth)
}

// WorkspaceProjectReferences ...
func WorkspaceProjectReferences(workspace string) ([]string, error) {
	projects := []string{}

	workspaceDir := filepath.Dir(workspace)

	xcworkspacedataPth := path.Join(workspace, "contents.xcworkspacedata")
	if exist, err := pathutil.IsPathExists(xcworkspacedataPth); err != nil {
		return []string{}, err
	} else if !exist {
		return []string{}, fmt.Errorf("contents.xcworkspacedata does not exist at: %s", xcworkspacedataPth)
	}

	xcworkspacedataStr, err := fileutil.ReadStringFromFile(xcworkspacedataPth)
	if err != nil {
		return []string{}, err
	}

	xcworkspacedataLines := strings.Split(xcworkspacedataStr, "\n")
	fileRefStart := false
	regexp := regexp.MustCompile(`location = "(.+):(.+).xcodeproj"`)

	for _, line := range xcworkspacedataLines {
		if strings.Contains(line, "<FileRef") {
			fileRefStart = true
			continue
		}

		if fileRefStart {
			fileRefStart = false
			matches := regexp.FindStringSubmatch(line)
			if len(matches) == 3 {
				projectName := matches[2]
				project := filepath.Join(workspaceDir, projectName+".xcodeproj")
				projects = append(projects, project)
			}
		}
	}

	sort.Strings(projects)

	return projects, nil
}

// WorkspaceSharedSchemes ...
func WorkspaceSharedSchemes(workspacePth string) ([]SchemeModel, error) {
	workspaceSharedSchemes, err := sharedSchemes(workspacePth)
	if err != nil {
		return []SchemeModel{}, err
	}

	projects, err := WorkspaceProjectReferences(workspacePth)
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		projectSharedSchemes, err := sharedSchemes(project)
		if err != nil {
			return []SchemeModel{}, err
		}

		for _, projectSharedScheme := range projectSharedSchemes {
			for _, workspaceSharedScheme := range workspaceSharedSchemes {
				if workspaceSharedScheme.Name == projectSharedScheme.Name {
					continue
				}
			}

			workspaceSharedSchemes = append(workspaceSharedSchemes, projectSharedScheme)
		}
	}

	return workspaceSharedSchemes, nil
}

func buildableReferenceIDs(schemePth string) ([]string, error) {
	scheme, err := xcscheme.Open(schemePth)
	if err != nil {
		return []string{}, err
	}

	var ids []string
	for _, entry := range scheme.BuildAction.BuildActionEntries {
		ids = append(ids, entry.BuildableReference.BlueprintIdentifier)
	}

	return ids, nil
}
