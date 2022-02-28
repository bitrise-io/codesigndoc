package ios

import (
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcodeproj"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcscheme"
)

func schemeHasAppClipTarget(project xcodeproj.XcodeProj, scheme xcscheme.Scheme) bool {
	for _, entry := range scheme.BuildAction.BuildActionEntries {
		target, found := project.Proj.Target(entry.BuildableReference.BlueprintIdentifier)
		if !found {
			log.TDebugf("no target found for blueprint ID (%s) project (%s)", entry.BuildableReference.BlueprintIdentifier, project.Path)
			continue
		}

		if target.CanExportAppClip() {
			return true
		}
	}

	return false
}

func shouldAppendExportAppClipStep(hasAppClip bool, exportMethod string) bool {
	return hasAppClip &&
		(exportMethod == "development" || exportMethod == "ad-hoc")
}

func appendExportAppClipStep(configBuilder *models.ConfigBuilderModel, workflowID models.WorkflowID) {
	exportXCArchiveStepInputModels := []envmanModels.EnvironmentItemModel{
		{ProjectPathInputKey: "$" + ProjectPathInputEnvKey},
		{SchemeInputKey: "$" + SchemeInputEnvKey},
		{ExportXCArchiveProductInputKey: ExportXCArchiveProductInputAppClipValue},
		{DistributionMethodInputKey: "$" + DistributionMethodEnvKey},
		{AutomaticCodeSigningKey: AutomaticCodeSigningValue},
	}
	configBuilder.AppendStepListItemsTo(workflowID, steps.ExportXCArchiveStepListItem(exportXCArchiveStepInputModels...))
}
