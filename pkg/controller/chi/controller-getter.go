// Copyright 2019 Altinity Ltd and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package chi

import (
	"fmt"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sLabels "k8s.io/apimachinery/pkg/labels"

	log "github.com/altinity/clickhouse-operator/pkg/announcer"
	api "github.com/altinity/clickhouse-operator/pkg/apis/clickhouse.altinity.com/v1"
	"github.com/altinity/clickhouse-operator/pkg/controller"
	commonLabeler "github.com/altinity/clickhouse-operator/pkg/model/common/tags/labeler"
)

// getConfigMap gets ConfigMap either by namespaced name or by labels
// TODO review byNameOnly params
func (c *Controller) getConfigMap(meta meta.Object, byNameOnly bool) (*core.ConfigMap, error) {
	get := c.configMapLister.ConfigMaps(meta.GetNamespace()).Get
	list := c.configMapLister.ConfigMaps(meta.GetNamespace()).List
	var objects []*core.ConfigMap

	// Check whether object with such name already exists
	obj, err := get(meta.GetName())

	if (obj != nil) && (err == nil) {
		// Object found by name
		return obj, nil
	}

	if !apiErrors.IsNotFound(err) {
		// Error, which is not related to "Object not found"
		return nil, err
	}

	// Object not found by name

	if byNameOnly {
		return nil, err
	}

	// Try to find by labels

	var selector k8sLabels.Selector
	if selector, err = commonLabeler.MakeSelectorFromObjectMeta(meta); err != nil {
		return nil, err
	}

	if objects, err = list(selector); err != nil {
		return nil, err
	}

	if len(objects) == 0 {
		return nil, apiErrors.NewNotFound(apps.Resource("ConfigMap"), meta.GetName())
	}

	if len(objects) == 1 {
		// Exactly one object found by labels
		return objects[0], nil
	}

	// Too much objects found by labels
	return nil, fmt.Errorf("too much objects found %d expecting 1", len(objects))
}

// getSecret gets secret
func (c *Controller) getSecret(secret *core.Secret) (*core.Secret, error) {
	return c.kubeClient.CoreV1().Secrets(secret.Namespace).Get(controller.NewContext(), secret.Name, controller.NewGetOptions())
}

// getPodsIPs gets all pod IPs
func (c *Controller) getPodsIPs(obj interface{}) (ips []string) {
	log.V(3).M(obj).F().S().Info("looking for pods IPs")
	defer log.V(3).M(obj).F().E().Info("looking for pods IPs")

	for _, pod := range c.kube.Pod().GetAll(obj) {
		if ip := pod.Status.PodIP; ip == "" {
			log.V(3).M(pod).F().Warning("Pod NO IP address found. Pod: %s/%s", pod.Namespace, pod.Name)
		} else {
			ips = append(ips, ip)
			log.V(3).M(pod).F().Info("Pod IP address found. Pod: %s/%s IP: %s", pod.Namespace, pod.Name, ip)
		}
	}
	return ips
}

// GetCHIByObjectMeta gets CHI by namespaced name
func (c *Controller) GetCHIByObjectMeta(meta meta.Object, isCHI bool) (*api.ClickHouseInstallation, error) {
	var chiName string
	var err error
	if isCHI {
		chiName = meta.GetName()
	} else {
		chiName, err = commonLabeler.GetCRNameFromObjectMeta(meta)
		if err != nil {
			return nil, fmt.Errorf("unable to find CHI by name: '%s'. More info: %v", meta.GetName(), err)
		}
	}

	return c.chopClient.ClickhouseV1().ClickHouseInstallations(meta.GetNamespace()).Get(controller.NewContext(), chiName, controller.NewGetOptions())
}