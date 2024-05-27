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
	"context"
	core "k8s.io/api/core/v1"
	kube "k8s.io/client-go/kubernetes"

	"github.com/altinity/clickhouse-operator/pkg/controller"
)

type KubeEventClickHouse struct {
	kubeClient kube.Interface
}

func NewKubeEventClickHouse(kubeClient kube.Interface) *KubeEventClickHouse {
	return &KubeEventClickHouse{
		kubeClient: kubeClient,
	}
}

func (c *KubeEventClickHouse) Create(ctx context.Context, event *core.Event) (*core.Event, error) {
	return c.kubeClient.CoreV1().Events(event.Namespace).Create(ctx, event, controller.NewCreateOptions())
}
