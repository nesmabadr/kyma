package v1alpha2

import (
	"testing"

	"github.com/stretchr/testify/require"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vrischmann/envconfig"

	"github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const v1alpha1ReplicasPresetLabel = "serverless.kyma-project.io/replicas-preset"

func TestSetDefaults(t *testing.T) {
	zero := int32(0)
	one := int32(1)
	two := int32(2)

	functionProfiles := `
{
	"python39": "L"
}
`
	functionReplicas := `
{
"S":{"min": 1,"max": 1}
}
`
	functionResources := `
{
"S":{"requestCpu": "25m","requestMemory": "32Mi","limitCpu": "50m","limitMemory": "64Mi"},
"M":{"requestCpu": "50m","requestMemory": "64Mi","limitCpu": "100m","limitMemory": "128Mi"},
"L":{"requestCpu": "100m","requestMemory": "128Mi","limitCpu": "200m","limitMemory": "256Mi"}
}
`

	LRuntimeResources := ResourceRequirementsBuilder{}.Limits("200m", "256Mi").Requests("100m", "128Mi").BuildCoreV1()

	MRuntimeResources := ResourceRequirementsBuilder{}.Limits("100m", "128Mi").Requests("50m", "64Mi").BuildCoreV1()

	buildResources := `
{
"slow":{"requestCpu": "350m","requestMemory": "350Mi","limitCpu": "700m","limitMemory": "700Mi"},
"normal":{"requestCpu": "700m","requestMemory": "700Mi","limitCpu": "1100m","limitMemory": "1100Mi"},
"fast":{"requestCpu": "1100m","requestMemory": "1100Mi", "limitCpu": "1800m","limitMemory": "1800Mi"}
}
`

	fastBuildResources := ResourceRequirementsBuilder{}.Limits("1800m", "1800Mi").Requests("1100m", "1100Mi").BuildCoreV1()

	for testName, testData := range map[string]struct {
		givenFunc    Function
		expectedFunc Function
	}{
		"Should do nothing": {
			givenFunc: Function{
				Spec: FunctionSpec{
					Runtime: NodeJs14,
					ResourceConfiguration: &ResourceConfiguration{
						Function: ResourceRequirementsBuilder{}.Limits("150m", "158Mi").Requests("90m", "84Mi").Build(),
						Build:    ResourceRequirementsBuilder{}.Limits("400m", "321Mi").Requests("374m", "300Mi").Build(),
					},
					ScaleConfig: &ScaleConfig{
						MinReplicas: &two,
						MaxReplicas: &two,
					},
				},
			},
			expectedFunc: Function{Spec: FunctionSpec{
				Runtime: NodeJs14,
				ResourceConfiguration: &ResourceConfiguration{
					Function: ResourceRequirementsBuilder{}.Limits("150m", "158Mi").Requests("90m", "84Mi").Build(),
					Build:    ResourceRequirementsBuilder{}.Limits("400m", "321Mi").Requests("374m", "300Mi").Build(),
				},
				ScaleConfig: &ScaleConfig{
					MinReplicas: &two,
					MaxReplicas: &two,
				},
			},
			},
		},
		"Should not change runtime type": {
			givenFunc: Function{
				Spec: FunctionSpec{
					Runtime: Python39,
					ResourceConfiguration: &ResourceConfiguration{
						Function: ResourceRequirementsBuilder{}.Limits("150m", "158Mi").Requests("90m", "84Mi").Build(),
						Build:    ResourceRequirementsBuilder{}.Limits("400m", "321Mi").Requests("374m", "300Mi").Build(),
					},
					ScaleConfig: &ScaleConfig{
						MinReplicas: &two,
						MaxReplicas: &two,
					},
				},
			},
			expectedFunc: Function{
				Spec: FunctionSpec{
					Runtime: Python39,
					ResourceConfiguration: &ResourceConfiguration{
						Function: ResourceRequirementsBuilder{}.Limits("150m", "158Mi").Requests("90m", "84Mi").Build(),
						Build:    ResourceRequirementsBuilder{}.Limits("400m", "321Mi").Requests("374m", "300Mi").Build(),
					},
					ScaleConfig: &ScaleConfig{
						MinReplicas: &two,
						MaxReplicas: &two,
					},
				},
			},
		},
		"Should not change empty runtime type to default": {
			givenFunc: Function{
				Spec: FunctionSpec{
					ResourceConfiguration: &ResourceConfiguration{
						Function: ResourceRequirementsBuilder{}.Limits("150m", "158Mi").Requests("90m", "84Mi").Build(),
						Build:    ResourceRequirementsBuilder{}.Limits("400m", "321Mi").Requests("374m", "300Mi").Build(),
					},
					ScaleConfig: &ScaleConfig{
						MinReplicas: &two,
						MaxReplicas: &two,
					},
				},
			},
			expectedFunc: Function{
				Spec: FunctionSpec{
					ResourceConfiguration: &ResourceConfiguration{
						Function: ResourceRequirementsBuilder{}.Limits("150m", "158Mi").Requests("90m", "84Mi").Build(),
						Build:    ResourceRequirementsBuilder{}.Limits("400m", "321Mi").Requests("374m", "300Mi").Build(),
					},
					ScaleConfig: &ScaleConfig{
						MinReplicas: &two,
						MaxReplicas: &two,
					},
				},
			},
		},
		"Should return default webhook": {
			givenFunc: Function{},
			expectedFunc: Function{
				Spec: FunctionSpec{
					ResourceConfiguration: &ResourceConfiguration{
						Function: ResourceRequirementsBuilder{}.Limits("100m", "128Mi").Requests("50m", "64Mi").Build(),
					},
					ScaleConfig: &ScaleConfig{
						MinReplicas: &one,
						MaxReplicas: &one,
					},
				},
			},
		},
		"Should fill missing fields": {
			givenFunc: Function{
				Spec: FunctionSpec{
					ResourceConfiguration: &ResourceConfiguration{
						Function: ResourceRequirementsBuilder{}.Requests("150m", "150Mi").Build(),
						Build:    ResourceRequirementsBuilder{}.Requests("1200m", "12000Mi").Build(),
					},
					Replicas: &two,
				},
			},
			expectedFunc: Function{
				Spec: FunctionSpec{
					ResourceConfiguration: &ResourceConfiguration{
						Function: ResourceRequirementsBuilder{}.Limits("150m", "150Mi").Requests("150m", "150Mi").Build(),
						Build:    ResourceRequirementsBuilder{}.Limits("1200m", "12000Mi").Requests("1200m", "12000Mi").Build(),
					},
					Replicas: &two,
				},
			},
		},
		"should consider maxReplicas and limits": {
			givenFunc: Function{
				Spec: FunctionSpec{
					ResourceConfiguration: &ResourceConfiguration{
						Function: ResourceRequirementsBuilder{}.Limits("15m", "15Mi").Build(),
						Build:    ResourceRequirementsBuilder{}.Limits("800m", "800Mi").Build(),
					},
					ScaleConfig: &ScaleConfig{
						MaxReplicas: &zero,
					},
				},
			},
			expectedFunc: Function{
				Spec: FunctionSpec{
					ResourceConfiguration: &ResourceConfiguration{
						Function: ResourceRequirementsBuilder{}.Limits("15m", "15Mi").Requests("15m", "15Mi").Build(),
						Build:    ResourceRequirementsBuilder{}.Limits("800m", "800Mi").Requests("700m", "700Mi").Build(),
					},
					ScaleConfig: &ScaleConfig{
						MinReplicas: &zero,
						MaxReplicas: &zero,
					},
				},
			},
		},
	} {
		t.Run(testName, func(t *testing.T) {
			// given
			g := gomega.NewWithT(t)
			config := &DefaultingConfig{}
			err := envconfig.Init(config)
			g.Expect(err).To(gomega.BeNil())

			functionReplicasPresets, err := ParseReplicasPresets(functionReplicas)
			g.Expect(err).To(gomega.BeNil())
			config.Function.Replicas.Presets = functionReplicasPresets

			functionResourcesPresets, err := ParseResourcePresets(functionResources)
			g.Expect(err).To(gomega.BeNil())
			config.Function.Resources.Presets = functionResourcesPresets

			buildResourcesPresets, err := ParseResourcePresets(buildResources)
			g.Expect(err).To(gomega.BeNil())
			config.BuildJob.Resources.Presets = buildResourcesPresets

			// when
			testData.givenFunc.Default(config)

			// then
			g.Expect(testData.givenFunc).To(gomega.Equal(testData.expectedFunc))
		})
	}

	testCases := map[string]struct {
		givenFunc    Function
		expectedFunc Function
	}{
		"Should properly set resources presets (using labels) - case with all fields": {
			givenFunc: Function{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						FunctionResourcesPresetLabel: "S",
						BuildResourcesPresetLabel:    "slow",
					},
				},
				Spec: FunctionSpec{
					Runtime: NodeJs14,
				},
			},
			expectedFunc: Function{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						FunctionResourcesPresetLabel: "S",
						BuildResourcesPresetLabel:    "slow",
					},
				}, Spec: FunctionSpec{
					Runtime: NodeJs14,
					ResourceConfiguration: &ResourceConfiguration{
						Function: ResourceRequirementsBuilder{}.Limits("50m", "64Mi").Requests("25m", "32Mi").Build(),
						Build:    ResourceRequirementsBuilder{}.Limits("700m", "700Mi").Requests("350m", "350Mi").Build(),
					},
					ScaleConfig: &ScaleConfig{
						MinReplicas: &one,
						MaxReplicas: &one,
					},
				},
			},
		},
		"Should properly merge resources presets (using labels) - case with all fields": {
			givenFunc: Function{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						FunctionResourcesPresetLabel: "S",
						BuildResourcesPresetLabel:    "slow",
					},
				},
				Spec: FunctionSpec{
					Runtime: NodeJs14,
					ResourceConfiguration: &ResourceConfiguration{
						Function: ResourceRequirementsBuilder{}.Requests("15m", "15Mi").Build(),
						Build:    ResourceRequirementsBuilder{}.Requests("250m", "250Mi").Build(),
					},
					ScaleConfig: &ScaleConfig{
						MinReplicas: &two,
					},
				},
			},
			expectedFunc: Function{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						FunctionResourcesPresetLabel: "S",
						BuildResourcesPresetLabel:    "slow",
					},
				}, Spec: FunctionSpec{
					Runtime: NodeJs14,
					ResourceConfiguration: &ResourceConfiguration{
						Function: ResourceRequirementsBuilder{}.Limits("50m", "64Mi").Requests("15m", "15Mi").Build(),
						Build:    ResourceRequirementsBuilder{}.Limits("700m", "700Mi").Requests("250m", "250Mi").Build(),
					},
					ScaleConfig: &ScaleConfig{
						MinReplicas: &two,
						MaxReplicas: &two,
					},
				},
			},
		},
		"Should properly merge resources presets (using labels) - case with concatenating missing values with default preset": {
			givenFunc: Function{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						FunctionResourcesPresetLabel: "L",
						BuildResourcesPresetLabel:    "fast",
					},
				},
				Spec: FunctionSpec{
					Runtime: NodeJs14,
				},
			},
			expectedFunc: Function{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						FunctionResourcesPresetLabel: "L",
						BuildResourcesPresetLabel:    "fast",
					},
				},
				Spec: FunctionSpec{
					Runtime: NodeJs14,
					ResourceConfiguration: &ResourceConfiguration{
						Function: &ResourceRequirements{
							Resources: &LRuntimeResources,
						},
						Build: &ResourceRequirements{
							Resources: &fastBuildResources,
						},
					},
					ScaleConfig: &ScaleConfig{
						MinReplicas: &one,
						MaxReplicas: &one,
					},
				},
			},
		},
		"Should set function profile to function presets M instead of default L value (using labels)": {
			givenFunc: Function{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						FunctionResourcesPresetLabel: "M",
					},
				},
				Spec: FunctionSpec{
					Runtime: Python39,
				},
			},
			expectedFunc: Function{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						FunctionResourcesPresetLabel: "M",
					},
				},
				Spec: FunctionSpec{
					Runtime: Python39,
					ResourceConfiguration: &ResourceConfiguration{
						Function: &ResourceRequirements{
							Resources: &MRuntimeResources,
						},
					},
					ScaleConfig: &ScaleConfig{
						MinReplicas: &one,
						MaxReplicas: &one,
					},
				}},
		},
		"Should properly merge resources presets (using labels) - case with missing buildResources Requests": {
			givenFunc: Function{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						FunctionResourcesPresetLabel: "S",
						BuildResourcesPresetLabel:    "slow",
					},
				},
				Spec: FunctionSpec{
					Runtime: NodeJs14,
					ResourceConfiguration: &ResourceConfiguration{
						Function: ResourceRequirementsBuilder{}.Requests("15m", "15Mi").Build(),
						Build:    ResourceRequirementsBuilder{}.Limits("700m", "700Mi").Build(),
					},
					ScaleConfig: &ScaleConfig{
						MinReplicas: &two,
					},
				},
			},
			expectedFunc: Function{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						FunctionResourcesPresetLabel: "S",
						BuildResourcesPresetLabel:    "slow",
					},
				}, Spec: FunctionSpec{
					Runtime: NodeJs14,
					ResourceConfiguration: &ResourceConfiguration{
						Function: ResourceRequirementsBuilder{}.Limits("50m", "64Mi").Requests("15m", "15Mi").Build(),
						Build:    ResourceRequirementsBuilder{}.Limits("700m", "700Mi").Requests("350m", "350Mi").Build(),
					},
					ScaleConfig: &ScaleConfig{
						MinReplicas: &two,
						MaxReplicas: &two,
					},
				},
			},
		},
		"Should ignore label replicas-preset": {
			givenFunc: Function{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						v1alpha1ReplicasPresetLabel: "XL",
					},
				},
				Spec: FunctionSpec{
					Runtime: NodeJs14,
				},
			},
			expectedFunc: Function{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						v1alpha1ReplicasPresetLabel: "XL",
					},
				}, Spec: FunctionSpec{
					Runtime: NodeJs14,
					ResourceConfiguration: &ResourceConfiguration{
						Function: ResourceRequirementsBuilder{}.Limits("100m", "128Mi").Requests("50m", "64Mi").Build(),
					},
					ScaleConfig: &ScaleConfig{
						MinReplicas: &one,
						MaxReplicas: &one,
					},
				},
			},
		},
	}

	for testName, testData := range testCases {
		t.Run(testName, func(t *testing.T) {
			// given
			g := gomega.NewWithT(t)
			config := &DefaultingConfig{}
			err := envconfig.Init(config)
			g.Expect(err).To(gomega.BeNil())

			functionReplicasPresets, err := ParseReplicasPresets(functionReplicas)
			g.Expect(err).To(gomega.BeNil())
			config.Function.Replicas.Presets = functionReplicasPresets

			functionResourcesPresets, err := ParseResourcePresets(functionResources)
			g.Expect(err).To(gomega.BeNil())
			config.Function.Resources.Presets = functionResourcesPresets

			buildResourcesPresets, err := ParseResourcePresets(buildResources)
			g.Expect(err).To(gomega.BeNil())
			config.BuildJob.Resources.Presets = buildResourcesPresets

			functionProfile, err := ParseRuntimePresets(functionProfiles)
			g.Expect(err).To(gomega.BeNil())
			config.Function.Resources.RuntimePresets = functionProfile
			// when
			testData.givenFunc.Default(config)

			// then
			//g.Expect(testData.givenFunc).To(gomega.Equal(testData.expectedFunc))
			require.EqualValues(t, testData.givenFunc, testData.expectedFunc)
		})
	}
}

type ResourceRequirementsBuilder struct {
	limitsCpu, limitsMemory, requestsCpu, requestsMemory string
}

func (b ResourceRequirementsBuilder) Limits(cpu, memory string) ResourceRequirementsBuilder {
	b.limitsCpu = cpu
	b.limitsMemory = memory
	return b
}

func (b ResourceRequirementsBuilder) Requests(cpu, memory string) ResourceRequirementsBuilder {
	b.requestsCpu = cpu
	b.requestsMemory = memory
	return b
}

func (b ResourceRequirementsBuilder) BuildCoreV1() corev1.ResourceRequirements {
	limits := corev1.ResourceList{}
	if b.limitsCpu != "" {
		limits[corev1.ResourceCPU] = resource.MustParse(b.limitsCpu)
	}
	if b.limitsMemory != "" {
		limits[corev1.ResourceMemory] = resource.MustParse(b.limitsMemory)
	}
	if len(limits) == 0 {
		limits = nil
	}
	requests := corev1.ResourceList{}
	if b.requestsCpu != "" {
		requests[corev1.ResourceCPU] = resource.MustParse(b.requestsCpu)
	}
	if b.requestsMemory != "" {
		requests[corev1.ResourceMemory] = resource.MustParse(b.requestsMemory)
	}
	if len(requests) == 0 {
		requests = nil
	}
	return corev1.ResourceRequirements{
		Limits:   limits,
		Requests: requests,
	}
}

func (b ResourceRequirementsBuilder) Build() *ResourceRequirements {
	res := b.BuildCoreV1()
	return &ResourceRequirements{
		Resources: &res,
	}
}