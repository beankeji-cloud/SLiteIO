// =======================================================================
// Copyright 2021 The LiteIO Authors.
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
// =======================================================================
// Modifications by The SLiteIO Authors on 2025:
// - Modification : support lvm thin volume

package priority

import (
	"context"

	v1 "lite.io/liteio/pkg/api/volume.antstor.alipay.com/v1"
	"lite.io/liteio/pkg/controller/manager/state"
	"k8s.io/klog/v2"
)

// PriorityByLeastResource is a PriorityFunc. Nodes with less free resource are more prefered.
func PriorityByLeastResource(ctx context.Context, n *state.Node, vol *v1.AntstorVolume) int {
	var score int
	// LeastResourcePoriotiy
	// the less free resource is remained, the larger the score is
	/*
		score = (allocated / total) * 100
	*/

	var total = float64(n.Pool.GetVgTotalBytes())
	var free = float64(n.Pool.GetVgFreeBytes())
	if total <= 0 {
		klog.Errorf("found StoragePool %s total space is invalid", n.Info.ID)
		return 0
	}
	if n.Pool.IsThin {
		total *= n.Pool.OverprovisionRatio
		free = float64(n.Pool.GetVgVirtualFreeBytes())
	}

	score = int((total - free) / total * 100)

	if score < 0 {
		score = 0
	}

	return score
}
