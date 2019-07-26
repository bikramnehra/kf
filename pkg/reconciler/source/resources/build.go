// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resources

import (
	"github.com/google/kf/pkg/apis/kf/v1alpha1"
	"github.com/knative/serving/pkg/resources"
	build "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/kmeta"
)

const (
	managedByLabel         = "app.kubernetes.io/managed-by"
	buildpackBuildTemplate = "buildpack"
	containerImageTemplate = "container"
)

// BuildName gets the name of a Build for a Source.
func BuildName(source *v1alpha1.Source) string {
	return source.Name
}

func makeContainerImageBuild(source *v1alpha1.Source) (*build.TaskRun, error) {
	buildName := BuildName(source)

	taskRef := &build.TaskRef{
		Name: "container",
		Kind: "ClusterTask",
	}

	buildOutput := []build.TaskResourceBinding{
		{
			Name: "OUTPUT_IMAGE",
			ResourceSpec: &build.PipelineResourceSpec{
				Type: build.PipelineResourceTypeImage,
				Params: []build.Param{
					{
						Name:  "url",
						Value: source.Spec.ContainerImage.Image,
					},
				},
			},
		},
	}

	return &build.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      buildName,
			Namespace: source.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*kmeta.NewControllerRef(source),
			},
			// Copy labels from the parent
			Labels: resources.UnionMaps(
				source.GetLabels(), map[string]string{
					managedByLabel: "kf",
				}),
		},
		Spec: build.TaskRunSpec{
			TaskRef: taskRef,
			Outputs: build.TaskRunOutputs{
				Resources: buildOutput,
			},
			ServiceAccount: source.Spec.ServiceAccount,
		},
	}, nil
}

func makeBuildpackBuild(source *v1alpha1.Source) (*build.TaskRun, error) {
	buildName := BuildName(source)
	appImageName := AppImageName(source)
	imageDestination := JoinRepositoryImage(source.Spec.BuildpackBuild.Registry, appImageName)

	taskRef := &build.TaskRef{
		Name: "buildpack",
		Kind: "ClusterTask",
	}

	buildInput := []build.TaskResourceBinding{
		{
			Name: "INPUT_IMAGE",
			ResourceSpec: &build.PipelineResourceSpec{
				Type: build.PipelineResourceTypeImage,
				Params: []build.Param{
					{
						Name:  "url",
						Value: source.Spec.ContainerImage.Image,
					},
				},
			},
		},
	}

	buildOutput := []build.TaskResourceBinding{
		{
			Name: "OUTPUT_IMAGE",
			ResourceSpec: &build.PipelineResourceSpec{
				Type: build.PipelineResourceTypeImage,
				Params: []build.Param{
					{
						Name:  "url",
						Value: imageDestination,
					},
				},
			},
		},
	}

	return &build.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      buildName,
			Namespace: source.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*kmeta.NewControllerRef(source),
			},
			// Copy labels from the parent
			Labels: resources.UnionMaps(
				source.GetLabels(), map[string]string{
					managedByLabel: "kf",
				}),
		},
		Spec: build.TaskRunSpec{
			TaskRef: taskRef,
			Inputs: build.TaskRunInputs{
				Resources: buildInput,
			},
			Outputs: build.TaskRunOutputs{
				Resources: buildOutput,
			},
			ServiceAccount: source.Spec.ServiceAccount,
		},
	}, nil
}

func makeObjectMeta(source *v1alpha1.Source) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      BuildName(source),
		Namespace: source.Namespace,
		OwnerReferences: []metav1.OwnerReference{
			*kmeta.NewControllerRef(source),
		},
		// Copy labels from the parent
		Labels: resources.UnionMaps(
			source.GetLabels(), map[string]string{
				managedByLabel: "kf",
			}),
	}
}

// MakeBuild creates a Build for a Source.
func MakeBuild(source *v1alpha1.Source) (*build.TaskRun, error) {
	if source.Spec.IsContainerBuild() {
		return makeContainerImageBuild(source)
	} else {
		return makeBuildpackBuild(source)
	}
}
