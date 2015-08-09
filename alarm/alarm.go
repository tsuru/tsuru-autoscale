// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package alarm

import (
	"errors"
	"fmt"
	stdlog "log"
	"sync"
	"time"

	"github.com/robertkrimen/otto"
	"github.com/tsuru/tsuru-autoscale/action"
	"github.com/tsuru/tsuru-autoscale/datasource"
	"github.com/tsuru/tsuru-autoscale/db"
	"github.com/tsuru/tsuru-autoscale/log"
	"github.com/tsuru/tsuru-autoscale/tsuru"
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
	Name       string            `json:"name"`
	Actions    []string          `json:"actions"`
	Expression string            `json:"expression"`
	Enabled    bool              `json:"enabled"`
	Wait       time.Duration     `json:"wait"`
	DataSource string            `json:"datasource"`
	Instance   string            `json:"instance"`
	Envs       map[string]string `json:"envs"`
}

func NewAlarm(a *Alarm) error {
	conn, err := db.Conn()
	if err != nil {
		return nil
	}
	defer conn.Close()
	return conn.Alarms().Insert(&a)
}

func runAutoScaleOnce() {
	logger().Print("checking alarms")
	alarms := []Alarm{}
	conn, err := db.Conn()
	if err != nil {
		return
	}
	defer conn.Close()
	err = conn.Alarms().Find(bson.M{"enabled": true}).All(&alarms)
	if err != nil {
		return
	}
	var wg sync.WaitGroup
	for _, alarm := range alarms {
		wg.Add(1)
		go func(alarm Alarm) {
			defer wg.Done()
			logger().Printf("checking %s alarm", alarm.Name)
			err := scaleIfNeeded(&alarm)
			if err != nil {
				logger().Print(err.Error())
			}
		}(alarm)
	}
	wg.Wait()
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
	logger().Printf("alarm %s - %s - check: %t", alarm.Name, alarm.Expression, check)
	if check {
		if wait, err := shouldWait(alarm); err != nil {
			logger().Printf("waiting for alarm %s", alarm.Name)
			return err
		} else if wait {
			return nil
		}
		for _, alarmName := range alarm.Actions {
			a, err := action.FindByName(alarmName)
			if err != nil {
				logger().Printf("alarm %s not found - error: %s", alarmName, err)
			} else {
				logger().Printf("executing alarm %s action %s", alarm.Name, a.Name)
				instance, err := tsuru.GetInstanceByName(alarm.Instance)
				if err != nil {
					logger().Printf("Error trying to get instance by name, auto scale aborted: %s", err)
					return err
				}
				if len(instance.Apps) < 1 {
					msg := "Error trying to get app instance, auto scale aborted."
					logger().Printf(msg)
					err = errors.New(msg)
					return err
				}
				appName := instance.Apps[0]
				evt, err := NewEvent(alarm, a)
				if err != nil {
					logger().Printf("Error trying to insert auto scale event, auto scale aborted: %s", err)
				}
				aErr := a.Do(appName, alarm.Envs)
				if aErr != nil {
					logger().Printf("Error executing action %s in the alarm %s - error: %s", a.Name, alarm.Name, aErr)
				} else {
					logger().Printf("alarm %s action %s executed", alarm.Name, a.Name)
				}
				err = evt.update(aErr)
				if err != nil {
					logger().Printf("Error trying to update auto scale event: %s", err)
				}
			}
		}
		return nil
	}
	return nil
}

func shouldWait(alarm *Alarm) (bool, error) {
	now := time.Now().UTC()
	lastEvent, err := lastScaleEvent(alarm)
	if err != nil && err != mgo.ErrNotFound {
		logger().Printf("error on get last event for alarm %s - not waiting - err: %s", alarm.Name, err)
		return false, err
	}
	if err != mgo.ErrNotFound && lastEvent.EndTime.IsZero() {
		logger().Printf("last event not finished yet for alarm %s - waiting", alarm.Name)
		return true, nil
	}
	diff := now.Sub(lastEvent.EndTime)
	if diff > alarm.Wait {
		logger().Printf("diff %d > %d form alarm %s - not waiting", diff, alarm.Wait, alarm.Name)
		return false, nil
	}
	logger().Printf("diff %d < %d form alarm %s - waiting", diff, alarm.Wait, alarm.Name)
	return true, nil
}

func Enable(alarm *Alarm) error {
	conn, err := db.Conn()
	if err != nil {
		return nil
	}
	defer conn.Close()
	return conn.Alarms().Update(bson.M{"name": alarm.Name}, bson.M{"$set": bson.M{"enabled": true}})
}

func Disable(alarm *Alarm) error {
	conn, err := db.Conn()
	if err != nil {
		return nil
	}
	defer conn.Close()
	return conn.Alarms().Update(bson.M{"name": alarm.Name}, bson.M{"$set": bson.M{"enabled": false}})
}

func (a *Alarm) Check() (bool, error) {
	logger().Printf("getting data for alarm %s", a.Name)
	ds, err := datasource.Get(a.DataSource)
	if err != nil {
		logger().Printf("error getting data for alarm %s - error: %s", a.Name, err.Error())
		return false, err
	}
	instance, err := tsuru.GetInstanceByName(a.Instance)
	if err != nil {
		logger().Printf("Error trying to get instance by name, auto scale aborted: %s", err)
		return false, err
	}
	if len(instance.Apps) < 0 {
		msg := "Error trying to get app instance."
		logger().Printf(msg)
		err = errors.New(msg)
		return false, err
	}
	appName := instance.Apps[0]
	data, err := ds.Get(appName)
	if err != nil {
		logger().Printf("error getting data for alarm %s - error: %s", a.Name, err.Error())
		return false, err
	}
	logger().Printf("data for alarm %s - %s", a.Name, data)
	vm := otto.New()
	vm.Run(fmt.Sprintf("var data=%s;", data))
	vm.Run(fmt.Sprintf("var expression=%s", a.Expression))
	expression, err := vm.Get("expression")
	if err != nil {
		logger().Printf("error executing expresion for alarm %s - error: %s", a.Name, err.Error())
		return false, err
	}
	check, err := expression.ToBoolean()
	if err != nil {
		return false, err
	}
	return check, nil
}

// ListAlarmsByToken lists alarms by token.
func ListAlarmsByToken(token string) ([]Alarm, error) {
	i, err := tsuru.FindServiceInstance(token)
	if err != nil {
		logger().Printf("error find service instance by token %s - error: %s", token, err.Error())
		return nil, err
	}
	instances := []string{}
	for _, instance := range i {
		instances = append(instances, instance.Name)
	}
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var alarms []Alarm
	err = conn.Alarms().Find(bson.M{"instance": bson.M{"$in": instances}}).All(&alarms)
	if err != nil {
		logger().Printf("error find alarms by instance #%v", instances)
		return nil, err
	}
	return alarms, nil
}

// ListAlarmsByInstance lists alarms by instance.
func ListAlarmsByInstance(instanceName string) ([]Alarm, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var alarms []Alarm
	err = conn.Alarms().Find(bson.M{"instance": instanceName}).All(&alarms)
	if err != nil {
		logger().Printf("error find alarms by instance %q", instanceName)
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

// RemoveAlarm removes an alarm.
func RemoveAlarm(a *Alarm) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Alarms().Remove(a)
}
