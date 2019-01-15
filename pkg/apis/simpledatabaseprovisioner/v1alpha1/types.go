/*
 * Copyright (c) 2019 Ecodia GmbH & Co. KG <opensource@ecodia.de>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SimpleDatabaseBinding is a specification for a SimpleDatabaseBinding resource
type SimpleDatabaseBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec SimpleDatabaseBindingSpec `json:"spec"`
}

// SimpleDatabaseBindingSpec is the spec for a SimpleDatabaseBinding resource
type SimpleDatabaseBindingSpec struct {
	InstanceName string `json:"instanceName"`
	SecretName   string `json:"secretName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SimpleDatabaseBindingList is a list of SimpleDatabaseBinding resources
type SimpleDatabaseBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []SimpleDatabaseBinding `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SimpleDatabaseInstance is a specification for a SimpleDatabaseInstance resource
type SimpleDatabaseInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec SimpleDatabaseInstanceSpec `json:"spec"`
}

// SimpleDatabaseInstanceSpec is the spec for a SimpleDatabaseInstance resource
type SimpleDatabaseInstanceSpec struct {
	DbmsServer   string `json:"dbmsServer"`
	DatabaseName string `json:"databaseName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SimpleDatabaseInstanceList is a list of SimpleDatabaseInstance resources
type SimpleDatabaseInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []SimpleDatabaseInstance `json:"items"`
}
