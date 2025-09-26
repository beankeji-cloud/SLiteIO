﻿// =======================================================================
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
// - Modification : fix aio bdev expantion bug; support mount option: discard

package spdk

import (
	"fmt"

	"lite.io/liteio/pkg/spdk/jsonrpc/client"
	"lite.io/liteio/pkg/util/misc"
	"k8s.io/klog/v2"
)

type AioBdevCreateRequest struct {
	BdevName  string
	DevPath   string
	BlockSize int
}

type AioBdevDeleteRequest struct {
	BdevName string
}

type AioBdevResizeRequest struct {
	BdevName string
	//TargetSize uint64
}

type AioServiceIface interface {
	CreateAioBdev(req AioBdevCreateRequest) (err error)
	DeleteAioBdev(req AioBdevDeleteRequest) (err error)
	ResizeAioBdev(req AioBdevResizeRequest) (err error)
}

func (svc *SpdkService) CreateAioBdev(req AioBdevCreateRequest) (err error) {
	svc.cli, err = svc.client()
	if err != nil {
		klog.Error("spdk client is nil, try to reconnect spdk socket", err)
		return
	}

	klog.Infof("creating aio_bdev. req is %+v", req)
	var (
		bdevName  = req.BdevName
		devPath   = req.DevPath
		blockSize = req.BlockSize
		fallocate = true
		hasDev    bool
	)
	// verify devPath exists
	hasDev, err = misc.FileExists(devPath)
	if err != nil || !hasDev {
		err = fmt.Errorf("devPath %s not exists, %t, %+v", devPath, hasDev, err)
		return
	}

	list, err := svc.cli.BdevGetBdevs(client.BdevGetBdevsReq{BdevName: bdevName})
	if err != nil {
		// only return error when error msg is "No such device"
		if !IsNotFoundDeviceError(err) {
			klog.Error(err)
			return
		}
	}
	klog.Infof("devPath %s bdevName %s", devPath, bdevName)

	for _, item := range list {
		if item.Name == bdevName {
			klog.Infof("devpath %s bdev %s already exists", devPath, bdevName)
			return
		}
	}

	_, err = svc.cli.BdevAioCreate(client.BdevAioCreateReq{
		BdevName:  bdevName,
		FileName:  devPath,
		BlockSize: blockSize,
		Fallocate: fallocate,
	})
	if err != nil {
		klog.Error(err)
		return
	}

	klog.Infof("created bdev %s", req.BdevName)
	return
}

func (svc *SpdkService) DeleteAioBdev(req AioBdevDeleteRequest) (err error) {
	svc.cli, err = svc.client()
	if err != nil {
		klog.Error("spdk client is nil, try to reconnect spdk socket", err)
		return
	}
	// delete aio bdev
	var bdevName = req.BdevName
	var foundBdev bool
	list, err := svc.cli.BdevGetBdevs(client.BdevGetBdevsReq{BdevName: bdevName})
	if err != nil {
		if !IsNotFoundDeviceError(err) {
			return
		}
	}

	if len(list) == 0 {
		klog.Infof("aio_bdev %+s is already deleted", bdevName)
		return nil
	}

	for _, item := range list {
		if item.Name == bdevName {
			foundBdev = true
			klog.Infof("found bdev %s to delete", bdevName)
			break
		}
	}

	if foundBdev {
		result, errRpc := svc.cli.BdevAioDelete(client.BdevAioDeleteReq{Name: bdevName})
		if errRpc != nil || !result {
			err = fmt.Errorf("delete AioBdev %s failed: %t, %+v", bdevName, result, errRpc)
			klog.Error(err)
			return
		}
	} else {
		klog.Infof("not found bdev %s, so consider it deleted", bdevName)
	}

	return
}

func (svc *SpdkService) ResizeAioBdev(req AioBdevResizeRequest) (err error) {
	svc.cli, err = svc.client()
	if err != nil {
		klog.Error("spdk client is nil, try to reconnect spdk socket", err)
		return
	}

	// check aio bdev
	_, err = svc.cli.BdevGetBdevs(client.BdevGetBdevsReq{BdevName: req.BdevName})
	if err != nil {
		klog.Error(err)
		return
	}

	var result bool
	result, err = svc.cli.BdevAioResize(client.BdevAioResizeReq{
		Name: req.BdevName,
		//Size: req.TargetSize,
	})
	if err != nil || !result {
		err = fmt.Errorf("resize AioBdev %s failed: %t, %+v", req.BdevName, result, err)
		klog.Error(err)
		return
	}

	return
}
