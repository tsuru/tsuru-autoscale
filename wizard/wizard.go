// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wizard

import (
	"fmt"
	"time"

	"github.com/tsuru/tsuru-autoscale/alarm"
)

type autoscale struct {
	name      string
	scaleUp   scaleAction
	scaleDown scaleAction
	minUnits  int
}

type scaleAction struct {
	metric   string
	operator string
	value    string
	step     string
	wait     time.Duration
}

func New(a autoscale) error {
	err := newScaleAction(a.scaleUp, "scale_up", a.name)
	if err != nil {
		return err
	}
	err = newScaleAction(a.scaleDown, "scale_down", a.name)
	if err != nil {
		return err
	}
	return nil
}

func enableScaleDown(instanceName string, minUnits int) error {
	a := alarm.Alarm{
		Name:       fmt.Sprintf("enable_scale_down_%s", instanceName),
		Expression: fmt.Sprintf("units > %d", minUnits),
		Enabled:    true,
		Wait:       15 * 1000 * 1000 * 1000,
		Actions:    []string{"enable_alarm"},
		Instance:   instanceName,
		Envs:       map[string]string{"alarm": fmt.Sprintf("scale_down_%s", instanceName)},
	}
	return alarm.NewAlarm(&a)
}

func disableScaleDown(instanceName string, minUnits int) error {
	a := alarm.Alarm{
		Name:       fmt.Sprintf("disable_scale_down_%s", instanceName),
		Expression: fmt.Sprintf("units <= %d", minUnits),
		Enabled:    true,
		Wait:       15 * 1000 * 1000 * 1000,
		Actions:    []string{"disable_alarm"},
		Instance:   instanceName,
		Envs:       map[string]string{"alarm": fmt.Sprintf("scale_down_%s", instanceName)},
	}
	return alarm.NewAlarm(&a)
}

func newScaleAction(action scaleAction, kind, instanceName string) error {
	a := alarm.Alarm{
		Name:       fmt.Sprintf("%s_%s", kind, instanceName),
		Expression: fmt.Sprintf("%s %s %s", action.metric, action.operator, action.value),
		Enabled:    true,
		Wait:       action.wait,
		Actions:    []string{kind},
		Instance:   instanceName,
		Envs:       map[string]string{"step": action.step},
	}
	return alarm.NewAlarm(&a)
}
