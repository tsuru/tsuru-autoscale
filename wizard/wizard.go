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
