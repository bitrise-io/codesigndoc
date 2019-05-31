package models

import bitriseModels "github.com/bitrise-io/bitrise/models"

type workflowBuilderModel struct {
	Steps       []bitriseModels.StepListItemModel
	Description string
}

func newDefaultWorkflowBuilder() *workflowBuilderModel {
	return &workflowBuilderModel{
		Steps: []bitriseModels.StepListItemModel{},
	}
}

func (builder *workflowBuilderModel) appendStepListItems(items ...bitriseModels.StepListItemModel) {
	builder.Steps = append(builder.Steps, items...)
}

func (builder *workflowBuilderModel) generate() bitriseModels.WorkflowModel {
	return bitriseModels.WorkflowModel{
		Steps:       builder.Steps,
		Description: builder.Description,
	}
}
