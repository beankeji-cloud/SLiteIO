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
// - Modification : support csi storage capacity tracking and Volume Health Monitoring

package driver

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
)

var (
	DefaultVolumeAccessModeType = []csi.VolumeCapability_AccessMode_Mode{
		csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
	}

	DefaultControllerServiceCapability = []csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_GET_CAPACITY,
		csi.ControllerServiceCapability_RPC_EXPAND_VOLUME,
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT,
		/*
			not support ControllerPublish/ControllerUnpublish, so the external-attacher will use trivialHandler
			to directly mark VolumeAttachment to attached status.
			Otherwise, external-attacher will use csiHandler. This handler will allocate a lot of memory to use pvLister, vaLister, csiNodeLister.

			csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
			csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT,
			csi.ControllerServiceCapability_RPC_CLONE_VOLUME,
		*/
	}

	DefaultNodeServiceCapability = []csi.NodeServiceCapability_RPC_Type{
		csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
		csi.NodeServiceCapability_RPC_EXPAND_VOLUME,
		csi.NodeServiceCapability_RPC_GET_VOLUME_STATS,
		csi.NodeServiceCapability_RPC_VOLUME_CONDITION,
	}

	DefaultPluginCapability = []*csi.PluginCapability{
		{
			Type: &csi.PluginCapability_Service_{
				Service: &csi.PluginCapability_Service{
					Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
				},
			},
		},
		//
		// https://kubernetes-csi.github.io/docs/topology.html#implementing-topology-in-your-csi-driver
		// if VOLUME_ACCESSIBILITY_CONSTRAINTS is enabled, the AccessibleTopology in CreateVolume Response will affect pod scheduling.
		// external-provisioner will set PV's nodeAffinity according to the value of AccessibleTopology. The scheduler will respect PV/s nodeAffinity.
		// Local PV is highly recommended to set nodeAffinity, so the PV will not be used by pod from other node.
		// CreateVolume 的返回中的 AccessibleTopology 会影响调度; 对于本地盘可以设置, 这样保证本地PV不会被其他节点pod使用;
		{
			Type: &csi.PluginCapability_Service_{
				Service: &csi.PluginCapability_Service{
					Type: csi.PluginCapability_Service_VOLUME_ACCESSIBILITY_CONSTRAINTS,
				},
			},
		},
		// LVM local PV supports online expansion. CSI ControllerExpandVolume and NodeExpandVolume must be implemented
		// 如果 node-attached volume 不支持在线扩容，那么需要声明这个 OFFLINE, 同时必须实现 ControllerExpandVolume 和 NodeExpandVolume
		{
			Type: &csi.PluginCapability_VolumeExpansion_{
				VolumeExpansion: &csi.PluginCapability_VolumeExpansion{
					Type: csi.PluginCapability_VolumeExpansion_ONLINE,
				},
			},
		},
	}
)

// NewVolumeCapabilityAccessMode creates CSI volume access mode object.
func NewVolumeCapabilityAccessMode(mode csi.VolumeCapability_AccessMode_Mode) *csi.VolumeCapability_AccessMode {
	return &csi.VolumeCapability_AccessMode{Mode: mode}
}

// NewControllerServiceCapability creates CSI controller capability object.
func NewControllerServiceCapability(cap csi.ControllerServiceCapability_RPC_Type) *csi.ControllerServiceCapability {
	return &csi.ControllerServiceCapability{
		Type: &csi.ControllerServiceCapability_Rpc{
			Rpc: &csi.ControllerServiceCapability_RPC{
				Type: cap,
			},
		},
	}
}

// NewNodeServiceCapability creates CSI node capability object.
func NewNodeServiceCapability(cap csi.NodeServiceCapability_RPC_Type) *csi.NodeServiceCapability {
	return &csi.NodeServiceCapability{
		Type: &csi.NodeServiceCapability_Rpc{
			Rpc: &csi.NodeServiceCapability_RPC{
				Type: cap,
			},
		},
	}
}
