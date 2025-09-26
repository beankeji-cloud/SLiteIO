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

package config

import v1 "lite.io/liteio/pkg/api/volume.antstor.alipay.com/v1"

const (
	AioBdevType  BdevType = "aioBdev"
	MemBdevType  BdevType = "memBdev"
	RaidBdevType BdevType = "raidBdev"

	DefaultLVMName   = "antstore-vg"
	DefaultLVSName   = "antstor_lvstore"
	DefaultRaid0Name = "antstor_raid0"

	DefaultMallocBdevName = "antstor_malloc"
	DefaultAioBdevName    = "antstor_aio"
)

var (
	DefaultLVM = StorageStack{
		Pooling: Pooling{
			Mode: v1.PoolModeKernelLVM,
			Name: DefaultLVMName,
		},
	}

	DefaultLVS = StorageStack{
		Pooling: Pooling{
			Mode: v1.PoolModeSpdkLVStore,
			Name: DefaultLVSName,
		},
		Bdev: &SpdkBdev{
			Type: RaidBdevType,
			Name: DefaultRaid0Name,
		},
	}
)

type BdevType string

type StorageStack struct {
	Pooling Pooling   `json:"pooling" yaml:"pooling"`
	PVs     []LvmPV   `json:"pvs,omitempty" yaml:"pvs"`
	Bdev    *SpdkBdev `json:"bdev,omitempty" yaml:"bdev"`
}

type Pooling struct {
	Mode v1.PoolMode `json:"mode" yaml:"mode"`
	Name string      `json:"name" yaml:"name"`
	IsThin       bool        `json:"isThin" yaml:"isThin"`
	ThinPoolName string      `json:"thinPoolName" yaml:"thinPoolName"`
	OverprovisionRatio float64     `json:"overprovisionRatio" yaml:"overprovisionRatio"`
}

type LvmPV struct {
	// DevicePath is device path of PV. if it is empty, create a loop device from a file
	DevicePath string `json:"devicePath" yaml:"devicePath"`
	// if DevicePath is empty, use Size to create a file
	Size uint64 `json:"size" yaml:"size"`
	// if not empty, create loop device from file
	FilePath         string `json:"filePath,omitempty" yaml:"filePath"`
	CreateIfNotExist bool   `json:"createIfNotExist,omitempty" yaml:"createIfNotExist"`
}

type SpdkBdev struct {
	Type BdevType `json:"type" yaml:"type"`
	Name string   `json:"name" yaml:"name"`
	// size in byte
	Size uint64 `json:"size" yaml:"size"`
	// for aioBdev
	FilePath         string `json:"filePath,omitempty" yaml:"filePath"`
	CreateIfNotExist bool   `json:"createIfNotExist,omitempty" yaml:"createIfNotExist"`
	// for vfio raidBdev
	VfioPCIeKeyword string `json:"vfioPCIeKeyword,omitempty" yaml:"vfioPCIeKeyword"`
}
