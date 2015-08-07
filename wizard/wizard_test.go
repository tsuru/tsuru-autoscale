// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wizard

import (
	"fmt"
	"os"
	"testing"

	"github.com/tsuru/tsuru-autoscale/alarm"
	"github.com/tsuru/tsuru-autoscale/db"
	"github.com/tsuru/tsuru/db/dbtest"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct {
	conn *db.Storage
}

func (s *S) SetUpSuite(c *check.C) {
	err := os.Setenv("MONGODB_DATABASE_NAME", "tsuru_autoscale_wizard")
	c.Assert(err, check.IsNil)
	s.conn, err = db.Conn()
	c.Assert(err, check.IsNil)
}

func (s *S) TearDownTest(c *check.C) {
	dbtest.ClearAllCollections(s.conn.Actions().Database)
}

func (s *S) TearDownSuite(c *check.C) {
	err := os.Unsetenv("MONGODB_DATABASE_NAME")
	c.Assert(err, check.IsNil)
}

var _ = check.Suite(&S{})

func (s *S) TestNewScale(c *check.C) {
	a := scaleAction{
		metric:   "cpu",
		operator: ">",
		step:     "1",
		value:    "10",
		wait:     50,
	}
	action := "scale_up"
	instanceName := "instanceName"
	scaleName := fmt.Sprintf("%s_%s", action, instanceName)
	err := newScaleAction(a, action, instanceName)
	c.Assert(err, check.IsNil)
	al, err := alarm.FindAlarmByName(scaleName)
	c.Assert(err, check.IsNil)
	c.Assert(al.Name, check.Equals, scaleName)
	c.Assert(al.Expression, check.Equals, fmt.Sprintf("%s %s %s", a.metric, a.operator, a.value))
	c.Assert(al.Envs, check.DeepEquals, map[string]string{"step": a.step})
	c.Assert(al.Enabled, check.Equals, true)
	c.Assert(al.Actions, check.DeepEquals, []string{action})
}

func (s *S) TestNew(c *check.C) {
	scaleUp := scaleAction{
		metric:   "cpu",
		operator: ">",
		step:     "1",
		value:    "10",
		wait:     50,
	}
	scaleDown := scaleAction{
		metric:   "cpu",
		operator: "<",
		step:     "1",
		value:    "2",
		wait:     50,
	}
	a := autoscale{
		name:      "test",
		scaleUp:   scaleUp,
		scaleDown: scaleDown,
	}
	err := New(a)
	c.Assert(err, check.IsNil)
	scaleName := "scale_up_test"
	al, err := alarm.FindAlarmByName(scaleName)
	c.Assert(err, check.IsNil)
	c.Assert(al.Name, check.Equals, scaleName)
	c.Assert(al.Expression, check.Equals, fmt.Sprintf("%s %s %s", scaleUp.metric, scaleUp.operator, scaleUp.value))
	c.Assert(al.Envs, check.DeepEquals, map[string]string{"step": scaleUp.step})
	c.Assert(al.Enabled, check.Equals, true)
	c.Assert(al.Actions, check.DeepEquals, []string{"scale_up"})
	scaleName = "scale_down_test"
	al, err = alarm.FindAlarmByName(scaleName)
	c.Assert(err, check.IsNil)
	c.Assert(al.Name, check.Equals, scaleName)
	c.Assert(al.Expression, check.Equals, fmt.Sprintf("%s %s %s", scaleDown.metric, scaleDown.operator, scaleDown.value))
	c.Assert(al.Envs, check.DeepEquals, map[string]string{"step": scaleDown.step})
	c.Assert(al.Enabled, check.Equals, true)
	c.Assert(al.Actions, check.DeepEquals, []string{"scale_down"})
}

func (s *S) TestEnableScaleDown(c *check.C) {
	minUnits := 2
	instanceName := "instanceName"
	err := enableScaleDown(instanceName, minUnits)
	c.Assert(err, check.IsNil)
	alarmName := fmt.Sprintf("enable_scale_down_%s", instanceName)
	al, err := alarm.FindAlarmByName(alarmName)
	c.Assert(err, check.IsNil)
	c.Assert(al.Name, check.Equals, alarmName)
	c.Assert(al.Expression, check.Equals, fmt.Sprintf("units > %d", minUnits))
	c.Assert(al.Envs, check.DeepEquals, map[string]string{"alarm": fmt.Sprintf("scale_down_%s", instanceName)})
	c.Assert(al.Enabled, check.Equals, true)
	c.Assert(al.Actions, check.DeepEquals, []string{"enable_alarm"})
}
