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
	if a.MinUnits <= 0 {
		a.MinUnits = 1
	}
	err := newScaleAction(a, "scale_up")
	if err != nil {
		logger().Error(err)
		return err
	}
	err = newScaleAction(a, "scale_down")
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

func newScaleAction(scaleConfig *AutoScale, kind string) error {
	var (
		name        string
		processName string
		action      ScaleAction
		expression  string
		datasources []string
	)
	if kind == "scale_up" {
		action = scaleConfig.ScaleUp
		datasources = []string{action.Metric}
	}
	if kind == "scale_down" {
		action = scaleConfig.ScaleDown
		expression = fmt.Sprintf(`!units.lock.Locked && units.units.map(function(unit){ if (unit.ProcessName === "{process}") {return 1} else {return 0}}).reduce(function(c, p) { return c + p }) > %d && `, scaleConfig.MinUnits)
		datasources = []string{action.Metric, "units"}
	}
	if scaleConfig.Process == "" {
		name = fmt.Sprintf("%s_%s", kind, scaleConfig.Name)
		processName = "web"
	} else {
		name = fmt.Sprintf("%s_%s_%s", kind, scaleConfig.Name, scaleConfig.Process)
		processName = scaleConfig.Process
	}
	aggregator := action.Aggregator
	if aggregator == "" {
		aggregator = "max"
	}
	expression += fmt.Sprintf("%s.aggregations.range.buckets[0].date.buckets[%s.aggregations.range.buckets[0].date.buckets.length - 1].%s.value %s %s", action.Metric, action.Metric, aggregator, action.Operator, action.Value)
	envs := map[string]string{
		"step":    action.Step,
		"process": processName,
	}
	a := alarm.Alarm{
		Name:        name,
		Expression:  expression,
		Enabled:     true,
		Wait:        action.Wait * time.Second,
		Actions:     []string{kind},
		Instance:    scaleConfig.Name,
		DataSources: datasources,
		Envs:        envs,
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
	var alarms []string
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
