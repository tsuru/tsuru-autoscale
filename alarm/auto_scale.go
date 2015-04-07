// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package alarm

import (
	"errors"
	"fmt"
	"time"

	"github.com/robertkrimen/otto"
	"github.com/tsuru/tsuru-autoscale/action"
	"github.com/tsuru/tsuru-autoscale/datasource"
	"github.com/tsuru/tsuru-autoscale/db"
	"github.com/tsuru/tsuru/log"
	"gopkg.in/mgo.v2"
)

func StartAutoScale() {
	go runAutoScale()
}

// Alarm represents the configuration for the auto scale.
type Alarm struct {
	Name       string              `json:"name"`
	Actions    []action.Action     `json:"actions"`
	Expression string              `json:"expression"`
	Enabled    bool                `json:"enabled"`
	Wait       time.Duration       `json:"wait"`
	DataSource datasource.Instance `json:"datasource"`
}

func NewAlarm(name, expression string, ds datasource.Instance) (*Alarm, error) {
	alarm := &Alarm{
		Name:       name,
		Expression: expression,
		Enabled:    true,
		DataSource: ds,
	}
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	err = conn.Alarms().Insert(&alarm)
	if err != nil {
		return nil, err
	}
	return alarm, nil
}

func runAutoScaleOnce() {
	alarms := []Alarm{}
	conn, err := db.Conn()
	if err != nil {
		return
	}
	defer conn.Close()
	err = conn.Alarms().Find(nil).All(&alarms)
	if err != nil {
		return
	}
	for _, alarm := range alarms {
		err := scaleIfNeeded(&alarm)
		if err != nil {
			log.Error(err.Error())
		}
	}
}

func runAutoScale() {
	for {
		runAutoScaleOnce()
		time.Sleep(30 * time.Second)
	}
}

func scaleIfNeeded(alarm *Alarm) error {
	if alarm == nil {
		return errors.New("alarm: alarm is not configured.")
	}
	check, err := alarm.Check()
	if err != nil {
		return err
	}
	if check {
		if wait, err := shouldWait(alarm, alarm.Wait); err != nil {
			return err
		} else if wait {
			return nil
		}
		evt, err := NewEvent(alarm, "increase")
		if err != nil {
			return fmt.Errorf("Error trying to insert auto scale event, auto scale aborted: %s", err.Error())
		}
		err = evt.update(nil)
		if err != nil {
			log.Errorf("Error trying to update auto scale event: %s", err.Error())
		}
		return nil
	}
	return nil
}

func shouldWait(alarm *Alarm, waitPeriod time.Duration) (bool, error) {
	now := time.Now().UTC()
	lastEvent, err := lastScaleEvent(alarm)
	if err != nil && err != mgo.ErrNotFound {
		return false, err
	}
	if err != mgo.ErrNotFound && lastEvent.EndTime.IsZero() {
		return true, nil
	}
	diff := now.Sub(lastEvent.EndTime)
	if diff > waitPeriod {
		return false, nil
	}
	return true, nil
}

func AutoScaleEnable(alarm *Alarm) error {
	alarm.Enabled = true
	return nil
}

func AutoScaleDisable(alarm *Alarm) error {
	alarm.Enabled = false
	return nil
}

func (a *Alarm) Check() (bool, error) {
	data, err := a.DataSource.Get()
	if err != nil {
		return false, err
	}
	vm := otto.New()
	vm.Run(fmt.Sprintf("var data=%s;", data))
	vm.Run(fmt.Sprintf("var expression=%s", a.Expression))
	expression, err := vm.Get("expression")
	if err != nil {
		return false, err
	}
	check, err := expression.ToBoolean()
	if err != nil {
		return false, err
	}
	return check, nil
}
