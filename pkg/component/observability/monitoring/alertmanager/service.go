// Copyright 2024 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package alertmanager

import (
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/utils/ptr"

	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	"github.com/gardener/gardener/pkg/component"
	gardenerutils "github.com/gardener/gardener/pkg/utils/gardener"
)

func (a *alertManager) service() *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      a.name(),
			Namespace: a.namespace,
			Labels:    a.getLabels(),
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: a.getLabels(),
			Ports: []corev1.ServicePort{{
				Name: PortNameMetrics,
				Port: port,
			}},
		},
	}

	networkPolicyPort := networkingv1.NetworkPolicyPort{
		Port:     ptr.To(intstr.FromInt32(port)),
		Protocol: ptr.To(corev1.ProtocolTCP),
	}

	switch a.values.ClusterType {
	case component.ClusterTypeSeed:
		utilruntime.Must(gardenerutils.InjectNetworkPolicyAnnotationsForGardenScrapeTargets(service, networkPolicyPort))
		utilruntime.Must(gardenerutils.InjectNetworkPolicyAnnotationsForSeedScrapeTargets(service, networkPolicyPort))
		utilruntime.Must(gardenerutils.InjectNetworkPolicyNamespaceSelectors(service, metav1.LabelSelector{MatchLabels: map[string]string{
			v1beta1constants.GardenRole: v1beta1constants.GardenRoleShoot,
		}}))

	case component.ClusterTypeShoot:
		utilruntime.Must(gardenerutils.InjectNetworkPolicyAnnotationsForScrapeTargets(service, networkPolicyPort))
	}

	return service
}
