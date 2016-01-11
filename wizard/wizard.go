// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wizard

import (
	"fmt"
	"time"

	"github.com/tsuru/tsuru-autoscale/alarm"
	"github.com/tsuru/tsuru-autoscale/db"
	"github.com/tsuru/tsuru-autoscale/log"
	"gopkg.in/mgo.v2/bson"
)

func logger() *log.Logger {
	return log.Log()
}

type AutoScale struct {
	Name      string      `json:"name"`
	ScaleUp   ScaleAction `json:"scaleUp"`
	ScaleDown ScaleAction `json:"scaleDown"`
	MinUnits  int         `json:"minUnits"`
	Process   string      `json:"process"`
}

type ScaleAction struct {
	Aggregator string        `json:"aggregator"`
	Metric     string        `json:"metric"`
	Operator   string        `json:"operator"`
	Value      string        `json:"value"`
	Step       string        `json:"step"`
	Wait       time.Duration `json:"wait"`
}

func New(a *AutoScale) error {
	err := newScaleAction(a.ScaleUp, "scale_up", a.Name, a.Process)
	if err != nil {
		logger().Error(err)
		return err
	}
	err = newScaleAction(a.ScaleDown, "scale_down", a.Name, a.Process)
	if err != nil {
		logger().Error(err)
		return err
	}
	err = enableScaleDown(a.Name, a.MinUnits, a.Process)
	if err != nil {
		logger().Error(err)
		return err
	}
	err = disableScaleDown(a.Name, a.MinUnits, a.Process)
	if err != nil {
		logger().Error(err)
		return err
	}
	conn, err := db.Conn()
	if err != nil {
		logger().Error(err)
		return nil
	}
	defer conn.Close()
	return conn.Wizard().Insert(&a)
}

func enableScaleDown(instanceName string, minUnits int, process string) error {
	var (
		name        string
		processName string
	)
	if process == "" {
		name = fmt.Sprintf("scale_down_%s", instanceName)
		processName = "web"
	} else {
		name = fmt.Sprintf("scale_down_%s_%s", instanceName, process)
		processName = process
	}
	a := alarm.Alarm{
		Name:        fmt.Sprintf("enable_scale_down_%s", instanceName),
		Expression:  fmt.Sprintf(`!units.lock.Locked && units.units.map(function(unit){ if (unit.ProcessName === "{process}") {return 1} else {return 0}}).reduce(function(c, p) { return c + p }) > %d`, minUnits),
		Enabled:     true,
		Wait:        15 * 1000 * 1000 * 1000,
		Actions:     []string{"enable_alarm"},
		Instance:    instanceName,
		DataSources: []string{"units"},
		Envs:        map[string]string{"alarm": name, "process": processName},
	}
	return alarm.NewAlarm(&a)
}

func disableScaleDown(instanceName string, minUnits int, process string) error {
	var (
		name        string
		processName string
	)
	if process == "" {
		name = fmt.Sprintf("scale_down_%s", instanceName)
		processName = "web"
	} else {
		name = fmt.Sprintf("scale_down_%s_%s", instanceName, process)
		processName = process
	}
	a := alarm.Alarm{
		Name:        fmt.Sprintf("disable_scale_down_%s", instanceName),
		Expression:  fmt.Sprintf(`!units.lock.Locked && units.units.map(function(unit){ if (unit.ProcessName === "{process}") {return 1} else {return 0}}).reduce(function(c, p) { return c + p }) <= %d`, minUnits),
		Enabled:     true,
		Wait:        15 * 1000 * 1000 * 1000,
		Actions:     []string{"disable_alarm"},
		Instance:    instanceName,
		DataSources: []string{"units"},
		Envs:        map[string]string{"alarm": name, "process": processName},
	}
	return alarm.NewAlarm(&a)
}

func newScaleAction(action ScaleAction, kind, instanceName, process string) error {
	var (
		name        string
		processName string
	)
	if process == "" {
		name = fmt.Sprintf("%s_%s", kind, instanceName)
		processName = "web"
	} else {
		name = fmt.Sprintf("%s_%s_%s", kind, instanceName, process)
		processName = process
	}
	aggregator := action.Aggregator
	if aggregator == "" {
		aggregator = "max"
	}
	a := alarm.Alarm{
		Name:        name,
		Expression:  fmt.Sprintf("%s.aggregations.range.buckets[0].date.buckets[%s.aggregations.range.buckets[0].date.buckets.length - 1].%s.value %s %s", action.Metric, action.Metric, aggregator, action.Operator, action.Value),
		Enabled:     true,
		Wait:        action.Wait * time.Second,
		Actions:     []string{kind},
		Instance:    instanceName,
		DataSources: []string{action.Metric},
		Envs: map[string]string{
			"step":    action.Step,
			"process": processName,
		},
	}
	return alarm.NewAlarm(&a)
}

// FindByName finds auto scale by name
func FindByName(name string) (*AutoScale, error) {
	conn, err := db.Conn()
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	defer conn.Close()
	var autoScale AutoScale
	err = conn.Wizard().Find(bson.M{"name": name}).One(&autoScale)
	if err != nil {
		logger().Error(err)
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
			logger().Error(err)
			return err
		}
		err = alarm.RemoveAlarm(al)
		if err != nil {
			logger().Error(err)
			return err
		}
	}
	return nil
}

// Remove removes an auto scale.
func Remove(a *AutoScale) error {
	err := removeAlarms(a)
	if err != nil {
		logger().Error(err)
		return err
	}
	conn, err := db.Conn()
	if err != nil {
		logger().Error(err)
		return err
	}
	defer conn.Close()
	return conn.Wizard().Remove(a)
}

func (a *AutoScale) Events() ([]alarm.Event, error) {
	conn, err := db.Conn()
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	defer conn.Close()
	var events []alarm.Event
	q := bson.M{"alarm.instance": a.Name, "alarm.actions": bson.M{"$in": []string{"scale_up", "scale_down"}}}
	err = conn.Events().Find(q).Sort("-starttime").Limit(200).All(&events)
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	return events, nil
}
