// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wizard

import (
	"fmt"
	"time"

	"github.com/tsuru/tsuru-autoscale/alarm"
	"github.com/tsuru/tsuru-autoscale/db"
)

type AutoScale struct {
	Name      string
	ScaleUp   scaleAction
	ScaleDown scaleAction
	MinUnits  int
}

type scaleAction struct {
	Metric   string
	Operator string
	Value    string
	Step     string
	Wait     time.Duration
}

func New(a *AutoScale) error {
	err := newScaleAction(a.ScaleUp, "scale_up", a.Name)
	if err != nil {
		return err
	}
	err = newScaleAction(a.ScaleDown, "scale_down", a.Name)
	if err != nil {
		return err
	}
	err = enableScaleDown(a.Name, a.MinUnits)
	if err != nil {
		return err
	}
	err = disableScaleDown(a.Name, a.MinUnits)
	if err != nil {
		return err
	}
	conn, err := db.Conn()
	if err != nil {
		return nil
	}
	defer conn.Close()
	return conn.Wizard().Insert(&a)
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
		Expression: fmt.Sprintf("%s %s %s", action.Metric, action.Operator, action.Value),
		Enabled:    true,
		Wait:       action.Wait,
		Actions:    []string{kind},
		Instance:   instanceName,
		Envs:       map[string]string{"step": action.Step},
	}
	return alarm.NewAlarm(&a)
}
