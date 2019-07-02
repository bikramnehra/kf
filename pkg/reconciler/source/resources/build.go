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
	"fmt"
	"time"

	"github.com/google/kf/pkg/apis/kf/v1alpha1"
	build "github.com/knative/build/pkg/apis/build/v1alpha1"
	"github.com/knative/serving/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/kmeta"
)

const (
	managedByLabel = "app.kubernetes.io/managed-by"
)

// BuildName gets the name of a Build for a Source.
func BuildName(source *v1alpha1.Source, timestamp int64) string {
	return fmt.Sprintf("%s-%d", source.Name, timestamp)
}

// AppImageName gets the image name for an application.
func AppImageName(namespace, appName string, timestamp int64) string {
	return fmt.Sprintf("app-%s-%s:%d", namespace, appName, timestamp)
}

// JoinRepositoryImage joins a repository and image name.
func JoinRepositoryImage(repository, imageName string) string {
	return fmt.Sprintf("%s/%s", repository, imageName)
}

// MakeBuild creates a Build for a Source.
func MakeBuild(source *v1alpha1.Source) (*build.Build, error) {

	buildSource := &build.SourceSpec{
		Custom: &corev1.Container{
			Image: source.Spec.BuildpackBuild.Source,
		},
	}

	timestamp := time.Now().UnixNano()
	appImageName := AppImageName(source.Namespace, source.Name, timestamp)
	imageDestination := JoinRepositoryImage("gcr.io/kf-source", appImageName)

	args := []build.ArgumentSpec{
		{
			Name:  "IMAGE",
			Value: imageDestination,
		},
		{
			Name:  "BUILDPACK",
			Value: source.Spec.BuildpackBuild.Buildpack,
		},
	}

	return &build.Build{
		ObjectMeta: metav1.ObjectMeta{
			Generation: 1,
			Name:       BuildName(source, timestamp),
			Namespace:  source.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*kmeta.NewControllerRef(source),
			},
			// Copy labels from the parent
			Labels: resources.UnionMaps(
				source.GetLabels(), map[string]string{
					managedByLabel: "kf",
					"kf-source":    source.Name,
				}),
		},
		Spec: build.BuildSpec{
			Source: buildSource,
			Template: &build.TemplateInstantiationSpec{
				Name:      "buildpack",
				Kind:      "ClusterBuildTemplate",
				Arguments: args,
			},
		},
	}, nil
}
