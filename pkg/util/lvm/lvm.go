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

package lvm

import "lite.io/liteio/pkg/util/osutil"

var (
	// default is cgo API; call EnableLvm2Cmd() to replace it with lvm2cmd implementation
	LvmUtil LvmIface = &cmd{
		jsonFormat: true,
		exec:       osutil.NewCommandExec(),
	}
)

type VG struct {
	Name      string
	UUID      string
	TotalByte uint64
	FreeByte  uint64
	PVCount   int
	// extends
	ExtendCount int
	ExtendSize  uint64
}

type LV struct {
	Name     string
	VGName   string
	DevPath  string
	SizeByte uint64
	// striped or linear or thin,pool
	LvLayout string
	// attributes
	LvAttr string
	// lv device status: "open" or ""
	LvDeviceOpen string
	// origin vol of snapshot
	Origin     string
	OriginUUID string
	OriginSize string
	DataPercent     string
	MetaDataPercent string
}

type LvOption struct {
	Size      uint64
	LogicSize string
}

type LvmIface interface {
	CreateVG(name string, pvs []string) (VG, error)
	CreatePV(pvs []string) error
	ListVG() ([]VG, error)
	ListLVInVG(vgName string) ([]LV, error)
	ListPV() ([]PV, error)
	CreateThinPool(vgName, poolName string) (err error)
	CreateThinLV(vgName, poolName, lvName string, sizeByte uint64) (vol LV, err error)
	CreateLinearLV(vgName, lvName string, opt LvOption) (vol LV, err error)
	CreateStripeLV(vgName, lvName string, sizeByte uint64) (vol LV, err error)
	RemoveLV(vgName, lvName string) (err error)
	RemoveVG(vgName string) (err error)
	RemovePVs(pvs []string) (err error)
	ExpandVolume(deltaBytes int64, targetVol string) (err error)

	CreateSnapshotLinear(vgName, snapName, originVol string, sizeByte uint64) (err error)
	CreateSnapshotStripe(vgName, snapName, originVol string, sizeByte uint64) (err error)
	MergeSnapshot(vgName, snapName string) (err error)
}
