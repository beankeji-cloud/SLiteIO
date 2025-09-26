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

package filter

import (
	"strconv"
	"strings"
	"sync"
)

const (
	ReasonPoolFreeSize      = "PoolFreeSize"
	ReasonSpdkUnhealthy     = "SpdkUnhealthy"
	ReasonRemoteVolMaxCount = "RemoteVolMaxCount"
	ReasonPositionNotMatch  = "PositionNotMatch"
	ReasonVolTypeNotMatch   = "VolTypeNotMatch"
	ReasonDataConflict      = "DataConflict"
	ReasonNodeAffinity      = "NodeAffinity"
	ReasonPoolAffinity      = "PoolAffinity"
	ReasonPoolUnschedulable = "PoolUnschedulable"
	ReasonReservationSize   = "ReservationTooSmall"
	ReasonReserveNotMatch   = "ReservationNotMatch"
	ReasonThinProvision     = "ThinProvision"

	NoStoragePoolAvailable = "NoStoragePoolAvailable"
	//
	CtxErrKey = "globalError"
)

type MergedError struct {
	// reason -> count
	reasons map[string]int
	lock    sync.Mutex
}

func NewMergedError() *MergedError {
	return &MergedError{
		reasons: map[string]int{},
	}
}

func (e *MergedError) Error() string {
	var s strings.Builder
	s.WriteString(NoStoragePoolAvailable + ": ")

	e.lock.Lock()
	defer e.lock.Unlock()
	for reason, cnt := range e.reasons {
		s.WriteString(reason + ": " + strconv.Itoa(cnt) + ", ")
	}

	return s.String()
}

func (e *MergedError) AddReason(reason string) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if cnt, has := e.reasons[reason]; has {
		e.reasons[reason] = cnt + 1
	} else {
		e.reasons[reason] = 1
	}
}

func IsNoStoragePoolAvailable(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), NoStoragePoolAvailable)
}
