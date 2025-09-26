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
// - Modification : support nvme connect parameters configurable

package driver

import (
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/klog/v2"
)

type CSIDriver struct {
	name               string
	version            string
	nodeID             string
	maxVolume          int64
	nvmeReconnectDelay int
	nvmeCtrlLossTMO    int
	volumeCap          []*csi.VolumeCapability_AccessMode
	controllerCap      []*csi.ControllerServiceCapability
	nodeCap            []*csi.NodeServiceCapability
	pluginCap          []*csi.PluginCapability
}

type NewCSIDriverOption struct {
	Name               string
	Version            string
	NodeID             string
	MaxVolume          int64
	NvmeReconnectDelay int
	NvmeCtrlLossTMO    int
	VolumeCap          []csi.VolumeCapability_AccessMode_Mode
	ControllerCap      []csi.ControllerServiceCapability_RPC_Type
	NodeCap            []csi.NodeServiceCapability_RPC_Type
	PluginCap          []*csi.PluginCapability
}

// NewCSIDriver create a CSI driver
func NewCSIDriver(opt NewCSIDriverOption) *CSIDriver {
	if opt.Name == "" {
		klog.Fatal("CSIDriverOption cannot be empty")
	}

	d := &CSIDriver{}
	d.name = opt.Name
	d.version = opt.Version
	// Setup Node Id
	d.nodeID = opt.NodeID
	// Setup max volume
	d.maxVolume = opt.MaxVolume
	d.nvmeReconnectDelay = opt.NvmeReconnectDelay
	d.nvmeCtrlLossTMO = opt.NvmeCtrlLossTMO
	// Setup cap
	d.addVolumeCapabilityAccessModes(opt.VolumeCap)
	d.addControllerServiceCapabilities(opt.ControllerCap)
	d.addNodeServiceCapabilities(opt.NodeCap)
	d.addPluginCapabilities(opt.PluginCap)

	return d
}

func (d *CSIDriver) addVolumeCapabilityAccessModes(vc []csi.VolumeCapability_AccessMode_Mode) {
	var vca []*csi.VolumeCapability_AccessMode
	for _, c := range vc {
		klog.V(4).Infof("Enabling volume access mode: %v", c.String())
		vca = append(vca, NewVolumeCapabilityAccessMode(c))
	}
	d.volumeCap = vca
}

func (d *CSIDriver) addControllerServiceCapabilities(cl []csi.ControllerServiceCapability_RPC_Type) {
	var csc []*csi.ControllerServiceCapability
	for _, c := range cl {
		klog.V(4).Infof("Enabling controller service capability: %v", c.String())
		csc = append(csc, NewControllerServiceCapability(c))
	}
	d.controllerCap = csc
}

func (d *CSIDriver) addNodeServiceCapabilities(nl []csi.NodeServiceCapability_RPC_Type) {
	var nsc []*csi.NodeServiceCapability
	for _, n := range nl {
		klog.V(4).Infof("Enabling node service capability: %v", n.String())
		nsc = append(nsc, NewNodeServiceCapability(n))
	}
	d.nodeCap = nsc
}

func (d *CSIDriver) addPluginCapabilities(cap []*csi.PluginCapability) {
	d.pluginCap = cap
}

func (d *CSIDriver) ValidateControllerServiceRequest(c csi.ControllerServiceCapability_RPC_Type) bool {
	if c == csi.ControllerServiceCapability_RPC_UNKNOWN {
		return true
	}

	for _, cap := range d.controllerCap {
		if c == cap.GetRpc().Type {
			return true
		}
	}
	return false
}

func (d *CSIDriver) ValidateNodeServiceRequest(c csi.NodeServiceCapability_RPC_Type) bool {
	if c == csi.NodeServiceCapability_RPC_UNKNOWN {
		return true
	}
	for _, cap := range d.nodeCap {
		if c == cap.GetRpc().Type {
			return true
		}
	}
	return false

}

func (d *CSIDriver) ValidateVolumeCapability(cap *csi.VolumeCapability) bool {
	return d.ValidateVolumeAccessMode(cap.GetAccessMode().GetMode())
}

func (d *CSIDriver) ValidateVolumeCapabilities(caps []*csi.VolumeCapability) bool {
	for _, cap := range caps {
		if !d.ValidateVolumeAccessMode(cap.GetAccessMode().GetMode()) {
			return false
		}
	}
	return true
}

func (d *CSIDriver) ValidateVolumeAccessMode(c csi.VolumeCapability_AccessMode_Mode) bool {
	for _, mode := range d.volumeCap {
		if c == mode.GetMode() {
			return true
		}
	}
	return false
}

func (d *CSIDriver) ValidatePluginCapabilityService(cap csi.PluginCapability_Service_Type) bool {
	for _, v := range d.GetPluginCapability() {
		if v.GetService() != nil && v.GetService().GetType() == cap {
			return true
		}
	}
	return false
}

func (d *CSIDriver) GetName() string {
	return d.name
}

func (d *CSIDriver) GetVersion() string {
	return d.version
}

func (d *CSIDriver) GetInstanceId() string {
	return d.nodeID
}

func (d *CSIDriver) GetMaxVolumePerNode() int64 {
	return d.maxVolume
}

func (d *CSIDriver) GetNvmeReconnectDelay() int {
	return d.nvmeReconnectDelay
}

func (d *CSIDriver) GetNvmeCtrlLossTMO() int {
	return d.nvmeCtrlLossTMO
}

func (d *CSIDriver) GetControllerCapability() []*csi.ControllerServiceCapability {
	return d.controllerCap
}

func (d *CSIDriver) GetNodeCapability() []*csi.NodeServiceCapability {
	return d.nodeCap
}

func (d *CSIDriver) GetPluginCapability() []*csi.PluginCapability {
	return d.pluginCap
}

func (d *CSIDriver) GetVolumeCapability() []*csi.VolumeCapability_AccessMode {
	return d.volumeCap
}

func (d *CSIDriver) GetTopologyZoneKey() string {
	return fmt.Sprintf("topology.%s/zone", d.GetName())
}

func (d *CSIDriver) GetTopologyInstanceTypeKey() string {
	return fmt.Sprintf("topology.%s/instance-type", d.GetName())
}
