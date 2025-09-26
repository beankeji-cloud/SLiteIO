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
	"fmt"

	v1 "lite.io/liteio/pkg/api/volume.antstor.alipay.com/v1"
	"lite.io/liteio/pkg/spdk"
	"lite.io/liteio/pkg/spdk/jsonrpc/client"
	"k8s.io/klog/v2"
)

type SpdkLvsPoolEngine struct {
	// LvsName is the name of LVS
	LvsName string
	// spdk service
	spdk spdk.SpdkServiceIface
}

func NewSpdkLvsPoolEngine(lvsName string, spdk spdk.SpdkServiceIface) (pe *SpdkLvsPoolEngine) {
	pe = &SpdkLvsPoolEngine{
		LvsName: lvsName,
		spdk:    spdk,
	}
	return
}

func (pe *SpdkLvsPoolEngine) PoolInfo(lvsName string) (info StaticInfo, err error) {
	var lvs client.LVStoreInfo
	lvs, err = pe.spdk.GetLVStore(lvsName)
	if err != nil {
		return
	}
	klog.Infof("found lvs, %+v", lvs)
	// assemble pool
	info.LVS = &v1.SpdkLVStore{
		Name:             lvs.Name,
		UUID:             lvs.UUID,
		BaseBdev:         lvs.BaseBdev,
		ClusterSize:      lvs.ClusterSize,
		TotalDataCluster: lvs.TotalDataClusters,
		BlockSize:        lvs.BlockSize,
		Bytes:            uint64(lvs.ClusterSize * lvs.TotalDataClusters),
	}

	return
}

func (pe *SpdkLvsPoolEngine) TotalAndFreeSize() (total, free, virtualFree uint64, dataPct, metadataPct float64, err error) {
	var lvs spdk.LVStoreInfo
	lvs, err = pe.spdk.GetLVStore(pe.LvsName)
	if err != nil {
		klog.Error(err)
	}
	total = uint64(lvs.ClusterSize * lvs.TotalDataClusters)
	free = uint64(lvs.ClusterSize * lvs.FreeClusters)
	virtualFree = free

	return
}

func (pe *SpdkLvsPoolEngine) CreateVolume(req CreateVolumeRequest) (resp CreateVolumeResponse, err error) {
	klog.Info("creating spdk lvol ", req)
	resp.UUID, err = pe.spdk.CreateLvol(spdk.CreateLvolReq{
		LVStore:  pe.LvsName,
		LvolName: req.VolName,
		SizeByte: int(req.SizeByte),
	})
	if err != nil {
		return
	}

	return
}

func (pe *SpdkLvsPoolEngine) DeleteVolume(volName string) (err error) {
	err = pe.spdk.DeleteLvol(spdk.DeleteLvolReq{
		LVStore:  pe.LvsName,
		LvolName: volName,
	})
	if err != nil {
		klog.Error(err)
		return
	}

	return
}

func (pe *SpdkLvsPoolEngine) GetVolume(volName string) (vol VolumeInfo, err error) {
	// get lvol bdev by name lvs/lvol
	var list []spdk.Bdev
	list, err = pe.spdk.BdevGetBdevs(spdk.BdevGetBdevsReq{
		BdevName: fmt.Sprintf("%s/%s", pe.LvsName, volName),
	})
	if err != nil {
		klog.Error(err)
		return
	}

	if len(list) > 0 {
		vol = VolumeInfo{
			Type: v1.VolumeTypeSpdkLVol,
			SpdkLvol: &SpdkLvolBdev{
				Lvol: v1.SpdkLvol{
					Name:    volName,
					LvsName: pe.LvsName,
					Thin:    false,
				},
				SizeByte: uint64(list[0].BlockSize * list[0].NumBlocks),
			},
		}
	}

	return
}

func (pe *SpdkLvsPoolEngine) CreateSnapshot(req CreateSnapshotRequest) (err error) {
	klog.Info("creating snapshot of Spdk lvol", req)
	_, err = pe.spdk.CreateLvolSnapshot(spdk.CreateLvolSnapReq{
		LvolFullName: req.OriginName,
		SnapName:     req.SnapshotName,
	})
	if err != nil {
		return
	}
	return
}

func (pe *SpdkLvsPoolEngine) RestoreSnapshot(snapshotName string) (err error) {
	err = fmt.Errorf("SPDK LVS not support RestoreSnapshot")
	return
}

func (pe *SpdkLvsPoolEngine) ExpandVolume(req ExpandVolumeRequest) (err error) {
	klog.Info("expanding SPDK lvol ", req)
	err = pe.spdk.ResizeLvol(spdk.ResizeLvolReq{
		LvolFullName: req.VolName,
		TargetSize:   req.TargetSize,
	})
	if err != nil {
		return
	}

	return
}
