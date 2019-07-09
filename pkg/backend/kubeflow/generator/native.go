// Copyright 2018 Caicloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	common "github.com/kubeflow/tf-operator/pkg/apis/common/v1"
	pytorchv1 "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1"
	tfv1 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/caicloud/ciao/pkg/types"
)

// Native is the type for native generator.
type Native struct {
	Namespace string
}

// New returns a new native generator.
func NewNative(namespace string) *Native {
	return &Native{
		Namespace: namespace,
	}
}

// GenerateTFJob generates a new TFJob.
func (n Native) GenerateTFJob(parameter *types.Parameter) (*tfv1.TFJob, error) {
	psCount := int32(parameter.PSCount)
	workerCount := int32(parameter.WorkerCount)
	cleanPodPolicy := common.CleanPodPolicy(parameter.CleanPolicy)

	psResource, err := parameter.Resource.PSLimits()
	if err != nil {
		return nil, err
	}
	workerResource, err := parameter.Resource.WorkerLimits()
	if err != nil {
		return nil, err
	}

	return &tfv1.TFJob{
		TypeMeta: metav1.TypeMeta{
			Kind: tfv1.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      parameter.GenerateName,
			Namespace: n.Namespace,
		},
		Spec: tfv1.TFJobSpec{
			CleanPodPolicy: &cleanPodPolicy,
			TFReplicaSpecs: map[tfv1.TFReplicaType]*common.ReplicaSpec{
				tfv1.TFReplicaTypePS: {
					Replicas: &psCount,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Name:  defaultContainerNameTF,
									Image: parameter.Image,
									Resources: v1.ResourceRequirements{
										Limits: psResource,
									},
								},
							},
						},
					},
				},
				tfv1.TFReplicaTypeWorker: {
					Replicas: &workerCount,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Name:  defaultContainerNameTF,
									Image: parameter.Image,
									Resources: v1.ResourceRequirements{
										Limits: workerResource,
									},
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

// GeneratePyTorchJob generates a new PyTorchJob.
func (n Native) GeneratePyTorchJob(parameter *types.Parameter) (*pytorchv1.PyTorchJob, error) {
	masterCount := int32(parameter.MasterCount)
	workerCount := int32(parameter.WorkerCount)
	cleanPodPolicy := common.CleanPodPolicy(parameter.CleanPolicy)

	masterResource, err := parameter.Resource.MasterLimits()
	if err != nil {
		return nil, err
	}
	workerResource, err := parameter.Resource.WorkerLimits()
	if err != nil {
		return nil, err
	}

	return &pytorchv1.PyTorchJob{
		TypeMeta: metav1.TypeMeta{
			Kind: pytorchv1.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      parameter.GenerateName,
			Namespace: n.Namespace,
		},
		Spec: pytorchv1.PyTorchJobSpec{
			CleanPodPolicy: &cleanPodPolicy,
			PyTorchReplicaSpecs: map[pytorchv1.PyTorchReplicaType]*common.ReplicaSpec{
				pytorchv1.PyTorchReplicaTypeMaster: {
					Replicas: &masterCount,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Name:  defaultContainerNamePyTorch,
									Image: parameter.Image,
									Resources: v1.ResourceRequirements{
										Limits: masterResource,
									},
								},
							},
						},
					},
				},
				pytorchv1.PyTorchReplicaTypeWorker: {
					Replicas: &workerCount,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Name:  defaultContainerNamePyTorch,
									Image: parameter.Image,
									Resources: v1.ResourceRequirements{
										Limits: workerResource,
									},
								},
							},
						},
					},
				},
			},
		},
	}, nil
}
