// =======================================================================
// Copyright 2025 The SLiteIO Authors.
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

package rpcserver

import (
	"fmt"
	"k8s.io/klog/v2"
	v1 "lite.io/liteio/pkg/api/volume.antstor.alipay.com/v1"
	"lite.io/liteio/pkg/csi/client"
	"sync"
	"time"
)

const (
	updateCapsTime = 30 * time.Second
)

type poolCap struct {
	totalSpace int64
	freeSpace  int64
	isThin     bool
}

type storagePoolCaps struct {
	sync.RWMutex
	pools map[string]*poolCap
	cli   client.AntstorClientIface
}

func newStoragePoolCaps(cli client.AntstorClientIface) *storagePoolCaps {
	s := &storagePoolCaps{
		pools: make(map[string]*poolCap),
		cli:   cli,
	}
	s.update()
	go func() {
		tick := time.Tick(updateCapsTime)
		for range tick {
			s.update()
		}
	}()
	return s
}

func (s *storagePoolCaps) update() {
	res, err := s.cli.ListStoragePool(v1.DefaultNamespace)
	if err != nil {
		klog.Infof("ListStoragePool err:%v", err)
		return
	}
	pools := make(map[string]*poolCap, len(res.Items))
	for _, sp := range res.Items {
		free := sp.GetVgFreeBytes()
		if sp.IsThin {
			vfree := sp.GetVgVirtualFreeBytes()
			if vfree < free {
				free = vfree
			}
		}
		pools[sp.Name] = &poolCap{
			totalSpace: sp.GetStorageBytes(),
			freeSpace:  free,
			isThin:     sp.IsThin,
		}
	}
	s.Lock()
	s.pools = pools
	s.Unlock()
}

func (s *storagePoolCaps) get(name string) (poolCap, error) {
	s.RLock()
	defer s.RUnlock()
	if p, ok := s.pools[name]; ok {
		return *p, nil
	} else {
		return poolCap{}, fmt.Errorf("storeage pool %s not exist", name)
	}
}

func (s *storagePoolCaps) getAll(isThin bool) int64 {
	s.RLock()
	defer s.RUnlock()
	var sum int64
	for _, p := range s.pools {
		if p.isThin == isThin {
			sum += p.freeSpace
		}
	}
	return sum
}
