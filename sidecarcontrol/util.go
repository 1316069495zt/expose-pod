/*
Copyright 2020 The Kruise Authors.

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

package sidecarcontrol

import (
	"regexp"

	corev1 "k8s.io/api/core/v1"
)

var (
	// SidecarIgnoredNamespaces specifies the namespaces where Pods won't get injected
	SidecarIgnoredNamespaces = []string{"kube-system", "kube-public"}
	// SubPathExprEnvReg format: $(ODD_NAME)、$(POD_NAME)...
	SubPathExprEnvReg, _ = regexp.Compile(`\$\(([-._a-zA-Z][-._a-zA-Z0-9]*)\)`)
)

// IsActivePod determines the pod whether need be injected and updated
func IsActivePod(pod *corev1.Pod) bool {
	for _, namespace := range SidecarIgnoredNamespaces {
		if pod.Namespace == namespace {
			return false
		}
	}
	if pod.ObjectMeta.GetDeletionTimestamp() != nil {
		return false
	}
	return true
}
