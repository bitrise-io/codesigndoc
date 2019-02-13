package xcodeproj

import (
	"bufio"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
)

// TargetModel ...
type TargetModel struct {
	Name      string
	HasXCTest bool
}

// PBXNativeTarget ...
type PBXNativeTarget struct {
	id           string
	isa          string
	dependencies []string
	name         string
	productPath  string
	productType  string
}

// PBXTargetDependency ...
type PBXTargetDependency struct {
	id     string
	isa    string
	target string
}

func parsePBXNativeTargets(pbxprojContent string) ([]PBXNativeTarget, error) {
	pbxNativeTargets := []PBXNativeTarget{}

	id := ""
	isa := ""
	dependencies := []string{}
	name := ""
	productPath := ""
	productType := ""

	beginPBXNativeTargetSectionPattern := `/* Begin PBXNativeTarget section */`
	endPBXNativeTargetSectionPattern := `/* End PBXNativeTarget section */`
	isPBXNativeTargetSection := false

	// BAAFFED019EE788800F3AC91 /* SampleAppWithCocoapods */ = {
	beginPBXNativeTargetRegexp := regexp.MustCompile(`\s*(?P<id>[A-Z0-9]+) /\* (?P<name>.*) \*/ = {`)
	endPBXNativeTargetPattern := `};`
	isPBXNativeTarget := false

	// isa = PBXNativeTarget;
	isaRegexp := regexp.MustCompile(`\s*isa = (?P<isa>.*);`)

	beginDependenciesPattern := `dependencies = (`
	dependencieRegexp := regexp.MustCompile(`\s*(?P<id>[A-Z0-9]+) /\* (?P<isa>.*) \*/,`)
	endDependenciesPattern := `);`
	isDependencies := false

	// name = SampleAppWithCocoapods;
	nameRegexp := regexp.MustCompile(`\s*name = (?P<name>.*);`)
	// productReference = BAAFFEED19EE788800F3AC91 /* SampleAppWithCocoapodsTests.xctest */;
	productReferenceRegexp := regexp.MustCompile(`\s*productReference = (?P<id>[A-Z0-9]+) /\* (?P<path>.*) \*/;`)
	// productType = "com.apple.product-type.bundle.unit-test";
	productTypeRegexp := regexp.MustCompile(`\s*productType = (?P<productType>.*);`)

	scanner := bufio.NewScanner(strings.NewReader(pbxprojContent))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == endPBXNativeTargetSectionPattern {
			break
		}

		if strings.TrimSpace(line) == beginPBXNativeTargetSectionPattern {
			isPBXNativeTargetSection = true
			continue
		}

		if !isPBXNativeTargetSection {
			continue
		}

		// PBXNativeTarget section

		if strings.TrimSpace(line) == endPBXNativeTargetPattern {
			pbxNativeTarget := PBXNativeTarget{
				id:           id,
				isa:          isa,
				dependencies: dependencies,
				name:         name,
				productPath:  productPath,
				productType:  productType,
			}
			pbxNativeTargets = append(pbxNativeTargets, pbxNativeTarget)

			id = ""
			isa = ""
			name = ""
			productPath = ""
			productType = ""
			dependencies = []string{}

			isPBXNativeTarget = false
			continue
		}

		if matches := beginPBXNativeTargetRegexp.FindStringSubmatch(line); len(matches) == 3 {
			id = matches[1]
			name = matches[2]

			isPBXNativeTarget = true
			continue
		}

		if !isPBXNativeTarget {
			continue
		}

		// PBXNativeTarget item

		if matches := isaRegexp.FindStringSubmatch(line); len(matches) == 2 {
			isa = strings.Trim(matches[1], `"`)
		}

		if matches := nameRegexp.FindStringSubmatch(line); len(matches) == 2 {
			name = strings.Trim(matches[1], `"`)
		}

		if matches := productTypeRegexp.FindStringSubmatch(line); len(matches) == 2 {
			productType = strings.Trim(matches[1], `"`)
		}

		if matches := productReferenceRegexp.FindStringSubmatch(line); len(matches) == 3 {
			// productId := strings.Trim(matches[1], `"`)
			productPath = strings.Trim(matches[2], `"`)
		}

		if isDependencies && strings.TrimSpace(line) == endDependenciesPattern {
			isDependencies = false
			continue
		}

		if strings.TrimSpace(line) == beginDependenciesPattern {
			isDependencies = true
			continue
		}

		if !isDependencies {
			continue
		}

		// dependencies
		if matches := dependencieRegexp.FindStringSubmatch(line); len(matches) == 3 {
			dependencieID := strings.Trim(matches[1], `"`)
			dependencieIsa := strings.Trim(matches[2], `"`)

			if dependencieIsa == "PBXTargetDependency" {
				dependencies = append(dependencies, dependencieID)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return []PBXNativeTarget{}, err
	}

	return pbxNativeTargets, nil
}

func parsePBXTargetDependencies(pbxprojContent string) ([]PBXTargetDependency, error) {
	pbxTargetDependencies := []PBXTargetDependency{}

	id := ""
	isa := ""
	target := ""

	beginPBXTargetDependencySectionPattern := `/* Begin PBXTargetDependency section */`
	endPBXTargetDependencySectionPattern := `/* End PBXTargetDependency section */`
	isPBXTargetDependencySection := false

	// BAAFFEEF19EE788800F3AC91 /* PBXTargetDependency */ = {
	beginPBXTargetDependencyRegexp := regexp.MustCompile(`\s*(?P<id>[A-Z0-9]+) /\* (?P<isa>.*) \*/ = {`)
	endPBXTargetDependencyPattern := `};`
	isPBXTargetDependency := false

	// isa = PBXTargetDependency;
	isaRegexp := regexp.MustCompile(`\s*isa = (?P<isa>.*);`)
	// target = BAAFFED019EE788800F3AC91 /* SampleAppWithCocoapods */;
	targetRegexp := regexp.MustCompile(`\s*target = (?P<id>[A-Z0-9]+) /\* (?P<name>.*) \*/;`)

	scanner := bufio.NewScanner(strings.NewReader(pbxprojContent))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == endPBXTargetDependencySectionPattern {
			break
		}

		if strings.TrimSpace(line) == beginPBXTargetDependencySectionPattern {
			isPBXTargetDependencySection = true
			continue
		}

		if !isPBXTargetDependencySection {
			continue
		}

		// PBXTargetDependency section

		if strings.TrimSpace(line) == endPBXTargetDependencyPattern {
			pbxTargetDependency := PBXTargetDependency{
				id:     id,
				isa:    isa,
				target: target,
			}
			pbxTargetDependencies = append(pbxTargetDependencies, pbxTargetDependency)

			id = ""
			isa = ""
			target = ""

			isPBXTargetDependency = false
			continue
		}

		if matches := beginPBXTargetDependencyRegexp.FindStringSubmatch(line); len(matches) == 3 {
			id = matches[1]
			isa = matches[2]

			isPBXTargetDependency = true
			continue
		}

		if !isPBXTargetDependency {
			continue
		}

		// PBXTargetDependency item

		if matches := isaRegexp.FindStringSubmatch(line); len(matches) == 2 {
			isa = strings.Trim(matches[1], `"`)
		}

		if matches := targetRegexp.FindStringSubmatch(line); len(matches) == 3 {
			targetID := strings.Trim(matches[1], `"`)
			// targetName := strings.Trim(matches[2], `"`)

			target = targetID
		}
	}

	return pbxTargetDependencies, nil
}

func targetDependencieWithID(dependencies []PBXTargetDependency, id string) (PBXTargetDependency, bool) {
	for _, dependencie := range dependencies {
		if dependencie.id == id {
			return dependencie, true
		}
	}
	return PBXTargetDependency{}, false
}

func targetWithID(targets []PBXNativeTarget, id string) (PBXNativeTarget, bool) {
	for _, target := range targets {
		if target.id == id {
			return target, true
		}
	}
	return PBXNativeTarget{}, false
}

func pbxprojContentTartgets(pbxprojContent string) ([]TargetModel, error) {
	targetMap := map[string]TargetModel{}

	nativeTargets, err := parsePBXNativeTargets(pbxprojContent)
	if err != nil {
		return []TargetModel{}, err
	}

	targetDependencies, err := parsePBXTargetDependencies(pbxprojContent)
	if err != nil {
		return []TargetModel{}, err
	}

	// Add targets which has test targets
	for _, target := range nativeTargets {
		if path.Ext(target.productPath) == ".xctest" {
			if len(target.dependencies) > 0 {
				for _, dependencieID := range target.dependencies {
					dependency, found := targetDependencieWithID(targetDependencies, dependencieID)
					if found {
						dependentTarget, found := targetWithID(nativeTargets, dependency.target)
						if found {
							targetMap[dependentTarget.name] = TargetModel{
								Name:      dependentTarget.name,
								HasXCTest: true,
							}
						}
					}
				}
			}
		}
	}

	// Add targets which has NO test targets
	for _, target := range nativeTargets {
		if path.Ext(target.productPath) != ".xctest" {
			_, found := targetMap[target.name]
			if !found {
				targetMap[target.name] = TargetModel{
					Name:      target.name,
					HasXCTest: false,
				}
			}
		}
	}

	targets := []TargetModel{}
	for _, target := range targetMap {
		targets = append(targets, target)
	}

	return targets, nil
}

// ProjectTargets ...
func ProjectTargets(projectPth string) ([]TargetModel, error) {
	pbxProjPth := filepath.Join(projectPth, "project.pbxproj")
	if exist, err := pathutil.IsPathExists(pbxProjPth); err != nil {
		return []TargetModel{}, err
	} else if !exist {
		return []TargetModel{}, fmt.Errorf("project.pbxproj does not exist at: %s", pbxProjPth)
	}

	content, err := fileutil.ReadStringFromFile(pbxProjPth)
	if err != nil {
		return []TargetModel{}, err
	}

	return pbxprojContentTartgets(content)
}

// WorkspaceTargets ...
func WorkspaceTargets(workspacePth string) ([]TargetModel, error) {
	projects, err := WorkspaceProjectReferences(workspacePth)
	if err != nil {
		return []TargetModel{}, err
	}

	targets := []TargetModel{}
	for _, project := range projects {
		projectTargets, err := ProjectTargets(project)
		if err != nil {
			return []TargetModel{}, err
		}

		targets = append(targets, projectTargets...)
	}

	return targets, nil
}
