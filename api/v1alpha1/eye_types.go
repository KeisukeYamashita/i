/*
Copyright 2019 KeisukeYamashita.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EyeSpec defines the desired state of Eye
type EyeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +kubebuilder:validation:Required
	Lifetime string `json:"lifetime"`

	SecretRef SecretRef `json:"secretRef,omitempty"`
}

// SecretRef ...
type SecretRef struct {
	Name string `json:"name"`
}

// EyeStatus defines the observed state of Eye
type EyeStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	NotExpired bool `json:"notExpired"`
}

// +kubebuilder:object:root=true

// Eye is the Schema for the eyes API
type Eye struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EyeSpec   `json:"spec,omitempty"`
	Status EyeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EyeList contains a list of Eye
type EyeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Eye `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Eye{}, &EyeList{})
}
