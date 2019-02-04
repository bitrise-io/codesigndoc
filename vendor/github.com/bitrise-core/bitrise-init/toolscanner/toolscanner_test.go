package toolscanner

import (
	"reflect"
	"testing"

	"github.com/bitrise-core/bitrise-init/models"
	bitriseModels "github.com/bitrise-io/bitrise/models"
	"github.com/google/go-cmp/cmp"
)

func TestAddProjectTypeToOptions(t *testing.T) {
	const detectedProjectType = "ios"
	type args struct {
		scannerOptionTree    models.OptionNode
		detectedProjectTypes []string
	}
	tests := []struct {
		name string
		args args
		want models.OptionNode
	}{
		{
			name: "1 project type",
			args: args{
				scannerOptionTree: models.OptionNode{
					Title:  "Working directory",
					EnvKey: "FASTLANE_WORK_DIR",
					ChildOptionMap: map[string]*models.OptionNode{
						"BitriseFastlaneSample": &models.OptionNode{
							Title:  "Fastlane lane",
							EnvKey: "FASTLANE_LANE",
							ChildOptionMap: map[string]*models.OptionNode{
								"ios test": &models.OptionNode{
									Config: "fastlane-config",
								},
							},
						},
					},
				},
				detectedProjectTypes: []string{detectedProjectType},
			},
			want: models.OptionNode{
				Title:  "Project type",
				EnvKey: ProjectTypeEnvKey,
				ChildOptionMap: map[string]*models.OptionNode{
					detectedProjectType: &models.OptionNode{
						Title:  "Working directory",
						EnvKey: "FASTLANE_WORK_DIR",
						ChildOptionMap: map[string]*models.OptionNode{
							"BitriseFastlaneSample": &models.OptionNode{
								Title:  "Fastlane lane",
								EnvKey: "FASTLANE_LANE",
								ChildOptionMap: map[string]*models.OptionNode{
									"ios test": &models.OptionNode{
										Config: "fastlane-config" + "_" + detectedProjectType,
									}}}}},
				},
			},
		},
		{
			name: "2 project types",
			args: args{
				scannerOptionTree: models.OptionNode{
					Title:  "Working directory",
					EnvKey: "FASTLANE_WORK_DIR",
					ChildOptionMap: map[string]*models.OptionNode{
						"BitriseFastlaneSample": &models.OptionNode{
							Title:  "Fastlane lane",
							EnvKey: "FASTLANE_LANE",
							ChildOptionMap: map[string]*models.OptionNode{
								"ios test": &models.OptionNode{
									Config: "fastlane-config",
								},
							},
						},
					},
				},
				detectedProjectTypes: []string{"ios", "android"},
			},
			want: models.OptionNode{
				Title:  "Project type",
				EnvKey: ProjectTypeEnvKey,
				ChildOptionMap: map[string]*models.OptionNode{
					"ios": &models.OptionNode{
						Title:  "Working directory",
						EnvKey: "FASTLANE_WORK_DIR",
						ChildOptionMap: map[string]*models.OptionNode{
							"BitriseFastlaneSample": &models.OptionNode{
								Title:  "Fastlane lane",
								EnvKey: "FASTLANE_LANE",
								ChildOptionMap: map[string]*models.OptionNode{
									"ios test": &models.OptionNode{
										Config: "fastlane-config" + "_" + "ios",
									},
								},
							},
						},
					},
					"android": &models.OptionNode{
						Title:  "Working directory",
						EnvKey: "FASTLANE_WORK_DIR",
						ChildOptionMap: map[string]*models.OptionNode{
							"BitriseFastlaneSample": &models.OptionNode{
								Title:  "Fastlane lane",
								EnvKey: "FASTLANE_LANE",
								ChildOptionMap: map[string]*models.OptionNode{
									"ios test": &models.OptionNode{
										Config: "fastlane-config" + "_" + "android",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddProjectTypeToOptions(tt.args.scannerOptionTree, tt.args.detectedProjectTypes); !reflect.DeepEqual(got.String(), tt.want.String()) {
				t.Errorf("AddProjectTypeToOptions() = %v, want %v", got, tt.want)
				t.Errorf("%s", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestAddProjectTypeToConfig(t *testing.T) {
	const title = "abcd"
	type args struct {
		configName           string
		config               bitriseModels.BitriseDataModel
		detectedProjectTypes []string
	}
	tests := []struct {
		name string
		args args
		want map[string]bitriseModels.BitriseDataModel
	}{
		{
			name: "2 project types",
			args: args{
				configName: "fastlane-config",
				config: bitriseModels.BitriseDataModel{
					Title:       title,
					ProjectType: "other",
				},
				detectedProjectTypes: []string{"ios", "android"},
			},
			want: map[string]bitriseModels.BitriseDataModel{
				"fastlane-config_ios": bitriseModels.BitriseDataModel{
					Title:       title,
					ProjectType: "ios",
				},
				"fastlane-config_android": bitriseModels.BitriseDataModel{
					Title:       title,
					ProjectType: "android",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddProjectTypeToConfig(tt.args.configName, tt.args.config, tt.args.detectedProjectTypes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddProjectTypeToConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
