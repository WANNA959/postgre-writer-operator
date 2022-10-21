/*
Copyright 2022.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PostgresWriterSpec defines the desired state of PostgresWriter
type PostgresWriterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Id of Student
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Type=integer
	Id int64 `json:"id,omitempty"`
	// id of Student
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Type=string
	Name string `json:"name,omitempty"`

	// age of Student
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Type=integer
	//+kubebuilder:validation:Minimum=0
	Age int32 `json:"age,omitempty"`

	// department of Student
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Type=string
	Department string `json:"department,omitempty"`
}

// PostgresWriterStatus defines the observed state of PostgresWriter
type PostgresWriterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// todo 添加get 信息
// +kubebuilder:printcolumn:name="Id",type="integer",JSONPath=".spec.id",description="id of item"
// +kubebuilder:printcolumn:name="Names",type="string",JSONPath=".spec.name",description="name of item"
// +kubebuilder:printcolumn:name="Ages",type="integer",JSONPath=".spec.age",description="age of item"
// +kubebuilder:printcolumn:name="Department",type="string",JSONPath=".spec.department",description="department of item"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="create time of crd"
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PostgresWriter is the Schema for the postgreswriters API
type PostgresWriter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PostgresWriterSpec   `json:"spec,omitempty"`
	Status PostgresWriterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PostgresWriterList contains a list of PostgresWriter
type PostgresWriterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PostgresWriter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PostgresWriter{}, &PostgresWriterList{})
}
