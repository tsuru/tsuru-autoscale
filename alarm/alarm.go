// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package alarm

import (
	"errors"
	"fmt"
	stdlog "log"
	"time"

	"github.com/robertkrimen/otto"
	"github.com/tsuru/tsuru-autoscale/action"
	"github.com/tsuru/tsuru-autoscale/datasource"
	"github.com/tsuru/tsuru-autoscale/db"
	"github.com/tsuru/tsuru-autoscale/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func StartAutoScale() {
	go runAutoScale()
}

func logger() *stdlog.Logger {
	return log.Logger()
}

// Alarm represents the configuration for the auto scale.
type Alarm struct {
	Name       string                `json:"name"`
	Actions    []action.Action       `json:"actions"`
	Expression string                `json:"expression"`
	Enabled    bool                  `json:"enabled"`
	Wait       time.Duration         `json:"wait"`
	DataSource datasource.DataSource `json:"datasource"`
	Instance   string                `json:"instance"`
}

func NewAlarm(a *Alarm) error {
	conn, err := db.Conn()
	if err != nil {
		return nil
	}
	defer conn.Close()
	err = conn.Alarms().Insert(&a)
	if err != nil {
		return nil
	}
	return nil
}

func runAutoScaleOnce() {
	logger().Print("checking alarms")
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
		logger().Printf("checking %s alarm", alarm.Name)
		err := scaleIfNeeded(&alarm)
		if err != nil {
			logger().Print(err.Error())
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
		logger().Printf("alarm %s check error: %s", alarm.Name, err.Error())
		return err
	}
	logger().Printf("alarm %s check: %t", alarm.Name, check)
	if check {
		if wait, err := shouldWait(alarm); err != nil {
			logger().Printf("waiting for alarm %s", alarm.Name)
			return err
		} else if wait {
			return nil
		}
		for _, a := range alarm.Actions {
			logger().Printf("executing alarm %s action %s", alarm.Name, a.Name)
			err := a.Do()
			if err != nil {
				logger().Printf("Error executing action %s in the alarm %s - error: ", a.Name, alarm.Name, err.Error())
			} else {
				logger().Printf("alarm %s action %s executed", alarm.Name, a.Name)
			}
		}
		evt, err := NewEvent(alarm)
		if err != nil {
			return fmt.Errorf("Error trying to insert auto scale event, auto scale aborted: %s", err.Error())
		}
		err = evt.update(nil)
		if err != nil {
			return fmt.Errorf("Error trying to update auto scale event: %s", err.Error())
		}
		return nil
	}
	return nil
}

func shouldWait(alarm *Alarm) (bool, error) {
	now := time.Now().UTC()
	lastEvent, err := lastScaleEvent(alarm)
	if err != nil && err != mgo.ErrNotFound {
		return false, err
	}
	if err != mgo.ErrNotFound && lastEvent.EndTime.IsZero() {
		return true, nil
	}
	diff := now.Sub(lastEvent.EndTime)
	if diff > alarm.Wait {
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
	logger().Printf("getting data for alarm %s", a.Name)
	data, err := a.DataSource.Get()
	if err != nil {
		logger().Printf("error getting data for alarm %s - error:", a.Name, err.Error())
		return false, err
	}
	logger().Printf("data for alarm %s", data)
	vm := otto.New()
	vm.Run(fmt.Sprintf("var data=%s;", data))
	vm.Run(fmt.Sprintf("var expression=%s", a.Expression))
	expression, err := vm.Get("expression")
	if err != nil {
		logger().Printf("error executing expresion for alarm %s - error:", a.Name, err.Error())
		return false, err
	}
	check, err := expression.ToBoolean()
	if err != nil {
		return false, err
	}
	return check, nil
}

// ListAlarms lists all alarms.
func ListAlarms() ([]Alarm, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var alarms []Alarm
	err = conn.Alarms().Find(nil).All(&alarms)
	if err != nil {
		return nil, err
	}
	return alarms, nil
}

// FindAlarmByName find alarm by name.
func FindAlarmByName(name string) (*Alarm, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var alarm Alarm
	err = conn.Alarms().Find(bson.M{"name": name}).One(&alarm)
	if err != nil {
		return nil, err
	}
	return &alarm, nil
}
