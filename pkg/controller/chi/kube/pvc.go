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

package kube

import (
	"context"
	core "k8s.io/api/core/v1"
	kube "k8s.io/client-go/kubernetes"

	"github.com/altinity/clickhouse-operator/pkg/controller"
)

type PVCClickHouse struct {
	kubeClient kube.Interface
}

func NewPVCClickHouse(kubeClient kube.Interface) *PVCClickHouse {
	return &PVCClickHouse{
		kubeClient: kubeClient,
	}
}

func (c *PVCClickHouse) Create(ctx context.Context, pvc *core.PersistentVolumeClaim) (*core.PersistentVolumeClaim, error) {
	return c.kubeClient.CoreV1().PersistentVolumeClaims(pvc.Namespace).Create(ctx, pvc, controller.NewCreateOptions())
}

func (c *PVCClickHouse) Get(ctx context.Context, namespace, name string) (*core.PersistentVolumeClaim, error) {
	return c.kubeClient.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, name, controller.NewGetOptions())
}

func (c *PVCClickHouse) Update(ctx context.Context, pvc *core.PersistentVolumeClaim) (*core.PersistentVolumeClaim, error) {
	return c.kubeClient.CoreV1().PersistentVolumeClaims(pvc.Namespace).Update(ctx, pvc, controller.NewUpdateOptions())
}

func (c *PVCClickHouse) Delete(ctx context.Context, namespace, name string) error {
	return c.kubeClient.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, name, controller.NewDeleteOptions())
}
