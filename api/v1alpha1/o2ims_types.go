/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// O2imsSpec defines the desired state of O2ims
type O2imsSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Url is the REST endpoint to call in case of reconcilation. Edit o2ims_types.go to remove/update
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Url string `json:"url,omitempty"`
}

// O2imsStatus defines the observed state of O2ims
type O2imsStatus struct {
	// Conditions store the status conditions of the O2ims instances
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// O2ims is the Schema for the o2ims API
type O2ims struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   O2imsSpec   `json:"spec,omitempty"`
	Status O2imsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// O2imsList contains a list of O2ims
type O2imsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []O2ims `json:"items"`
}

func init() {
	SchemeBuilder.Register(&O2ims{}, &O2imsList{})
}
