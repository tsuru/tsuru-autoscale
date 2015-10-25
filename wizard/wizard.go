// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wizard

import (
	"fmt"
	"time"

	"github.com/tsuru/tsuru-autoscale/alarm"
	"github.com/tsuru/tsuru-autoscale/db"
	"gopkg.in/mgo.v2/bson"
)

type AutoScale struct {
	Name      string      `json:"name"`
	ScaleUp   ScaleAction `json:"scaleUp"`
	ScaleDown ScaleAction `json:"scaleDown"`
	MinUnits  int         `json:"minUnits"`
	Process   string      `json:"process"`
}

type ScaleAction struct {
	Metric   string        `json:"metric"`
	Operator string        `json:"operator"`
	Value    string        `json:"value"`
	Step     string        `json:"step"`
	Wait     time.Duration `json:"wait"`
}

func New(a *AutoScale) error {
	err := newScaleAction(a.ScaleUp, "scale_up", a.Name, a.Process)
	if err != nil {
		return err
	}
	err = newScaleAction(a.ScaleDown, "scale_down", a.Name, a.Process)
	if err != nil {
		return err
	}
	err = enableScaleDown(a.Name, a.MinUnits, a.Process)
	if err != nil {
		return err
	}
	err = disableScaleDown(a.Name, a.MinUnits, a.Process)
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

func enableScaleDown(instanceName string, minUnits int, process string) error {
	var name string
	if process == "" {
		name = fmt.Sprintf("scale_down_%s", instanceName)
	} else {
		name = fmt.Sprintf("scale_down_%s_%s", instanceName, process)
	}
	a := alarm.Alarm{
		Name:       fmt.Sprintf("enable_scale_down_%s", instanceName),
		Expression: fmt.Sprintf("data.aggregations.range.buckets[0].date.buckets[0].unit.value > %d", minUnits),
		Enabled:    true,
		Wait:       15 * 1000 * 1000 * 1000,
		Actions:    []string{"enable_alarm"},
		Instance:   instanceName,
		DataSource: "units",
		Envs:       map[string]string{"alarm": name},
	}
	return alarm.NewAlarm(&a)
}

func disableScaleDown(instanceName string, minUnits int, process string) error {
	var name string
	if process == "" {
		name = fmt.Sprintf("scale_down_%s", instanceName)
	} else {
		name = fmt.Sprintf("scale_down_%s_%s", instanceName, process)
	}
	a := alarm.Alarm{
		Name:       fmt.Sprintf("disable_scale_down_%s", instanceName),
		Expression: fmt.Sprintf("data.aggregations.range.buckets[0].date.buckets[0].unit.value <= %d", minUnits),
		Enabled:    true,
		Wait:       15 * 1000 * 1000 * 1000,
		Actions:    []string{"disable_alarm"},
		Instance:   instanceName,
		DataSource: "units",
		Envs:       map[string]string{"alarm": name},
	}
	return alarm.NewAlarm(&a)
}

func newScaleAction(action ScaleAction, kind, instanceName, process string) error {
	var name string
	if process == "" {
		name = fmt.Sprintf("%s_%s", kind, instanceName)
	} else {
		name = fmt.Sprintf("%s_%s_%s", kind, instanceName, process)
	}
	a := alarm.Alarm{
		Name:       name,
		Expression: fmt.Sprintf("data.aggregations.range.buckets[0].date.buckets[0].max.value %s %s", action.Operator, action.Value),
		Enabled:    true,
		Wait:       action.Wait,
		Actions:    []string{kind},
		Instance:   instanceName,
		DataSource: action.Metric,
		Envs: map[string]string{
			"step":    action.Step,
			"process": process,
		},
	}
	return alarm.NewAlarm(&a)
}

// FindByName finds auto scale by name
func FindByName(name string) (*AutoScale, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var autoScale AutoScale
	err = conn.Wizard().Find(bson.M{"name": name}).One(&autoScale)
	if err != nil {
		return nil, err
	}
	return &autoScale, nil
}

func (a *AutoScale) alarms() []string {
	alarms := []string{
		fmt.Sprintf("enable_scale_down_%s", a.Name),
		fmt.Sprintf("disable_scale_down_%s", a.Name),
	}
	if a.Process == "" {
		alarms = append(alarms, fmt.Sprintf("scale_up_%s", a.Name))
		alarms = append(alarms, fmt.Sprintf("scale_down_%s", a.Name))
	} else {
		alarms = append(alarms, fmt.Sprintf("scale_up_%s_%s", a.Name, a.Process))
		alarms = append(alarms, fmt.Sprintf("scale_down_%s_%s", a.Name, a.Process))
	}
	return alarms
}

func removeAlarms(autoScale *AutoScale) error {
	for _, a := range autoScale.alarms() {
		al, err := alarm.FindAlarmByName(a)
		if err != nil {
			return err
		}
		err = alarm.RemoveAlarm(al)
		if err != nil {
			return err
		}
	}
	return nil
}

// Remove removes an auto scale.
func Remove(a *AutoScale) error {
	err := removeAlarms(a)
	if err != nil {
		return err
	}
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Wizard().Remove(a)
}
