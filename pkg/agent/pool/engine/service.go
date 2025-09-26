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

package engine

import (
	v1 "lite.io/liteio/pkg/api/volume.antstor.alipay.com/v1"
)

type VolumeServiceIface interface {
	CreateVolume(req CreateVolumeRequest) (resp CreateVolumeResponse, err error)
	DeleteVolume(volName string) (err error)
	GetVolume(volName string) (vol VolumeInfo, err error)
	CreateSnapshot(req CreateSnapshotRequest) (err error)
	RestoreSnapshot(snapshotName string) (err error)
	ExpandVolume(req ExpandVolumeRequest) (err error)
}

type PoolingInfoIface interface {
	PoolInfo(poolName string) (info StaticInfo, err error)
	TotalAndFreeSize() (total uint64, free, virtualFree uint64, dataPct, metadataPct float64, err error)
}

type PoolEngineIface interface {
	PoolingInfoIface
	VolumeServiceIface
}

type VolumeInfo struct {
	Type     v1.VolumeType
	LvmLV    *v1.KernelLVol
	SpdkLvol *SpdkLvolBdev
}

type SpdkLvolBdev struct {
	Lvol     v1.SpdkLvol
	SizeByte uint64
}

type StaticInfo struct {
	LVM *v1.KernelLVM
	LVS *v1.SpdkLVStore
}

type CreateVolumeRequest struct {
	// for LVM and SpdkLVS
	VolName  string
	SizeByte uint64
	// FsType to mkfs. Optional for LVM
	FsType string
	// LvLayout of lv to create. Optional for LVM
	LvLayout v1.LVLayout
}

type CreateVolumeResponse struct {
	// for SpdkLVS
	UUID string
	// for LVM
	DevPath string
}

type CreateSnapshotRequest struct {
	SnapshotName string
	OriginName   string
	SizeByte     uint64
}

type ExpandVolumeRequest struct {
	VolName    string
	TargetSize uint64
	OriginSize uint64
}
