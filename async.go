// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 Canonical Ltd
//
// SPDX-License-Identifier: Apache-2.0
//
package gxds

import (
	"fmt"

	"github.com/edgexfoundry/edgex-go/core/domain/models"
)

// processAsyncResults processes readings that are pushed from
// a DS implementation. Each is reading is optionally transformed
// before being pushed to Core Data.
func (s *Service) processAsyncResults() {
	for !s.stopped {
		readings := make([]models.Reading, 0, s.c.Device.MaxCmdOps)
		cr := <-s.asyncCh

		// get the device resource associated with the rsp.RO
		do := pc.getDeviceObjectByName(cr.DeviceName, cr.RO)

		_ = cr.TransformResult(do.Properties.Value)

		reading := cr.Reading(cr.DeviceName, do.Name)
		readings = append(readings, *reading)

		// push to Core Data
		event := &models.Event{Device: cr.DeviceName, Readings: readings}
		_, err := s.ec.Add(event)
		if err != nil {
			msg := fmt.Sprintf("internal error; failed to push event for dev: %s to CoreData: %s", cr.DeviceName, err)
			svc.lc.Error(msg)
		}
	}
}
