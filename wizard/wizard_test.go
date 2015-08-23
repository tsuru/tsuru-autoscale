// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wizard

import (
	"encoding/json"
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
		Metric:   "cpu",
		Operator: ">",
		Step:     "1",
		Value:    "10",
		Wait:     50,
	}
	action := "scale_up"
	instanceName := "instanceName"
	scaleName := fmt.Sprintf("%s_%s", action, instanceName)
	err := newScaleAction(a, action, instanceName)
	c.Assert(err, check.IsNil)
	al, err := alarm.FindAlarmByName(scaleName)
	c.Assert(err, check.IsNil)
	c.Assert(al.Name, check.Equals, scaleName)
	c.Assert(al.Expression, check.Equals, fmt.Sprintf("%s %s %s", a.Metric, a.Operator, a.Value))
	c.Assert(al.Envs, check.DeepEquals, map[string]string{"step": a.Step})
	c.Assert(al.Enabled, check.Equals, true)
	c.Assert(al.Actions, check.DeepEquals, []string{action})
}

func (s *S) TestNew(c *check.C) {
	scaleUp := scaleAction{
		Metric:   "cpu",
		Operator: ">",
		Step:     "1",
		Value:    "10",
		Wait:     50,
	}
	scaleDown := scaleAction{
		Metric:   "cpu",
		Operator: "<",
		Step:     "1",
		Value:    "2",
		Wait:     50,
	}
	a := AutoScale{
		Name:      "test",
		ScaleUp:   scaleUp,
		ScaleDown: scaleDown,
	}
	err := New(&a)
	c.Assert(err, check.IsNil)
	scaleName := "scale_up_test"
	al, err := alarm.FindAlarmByName(scaleName)
	c.Assert(err, check.IsNil)
	c.Assert(al.Name, check.Equals, scaleName)
	c.Assert(al.Expression, check.Equals, fmt.Sprintf("%s %s %s", scaleUp.Metric, scaleUp.Operator, scaleUp.Value))
	c.Assert(al.Envs, check.DeepEquals, map[string]string{"step": scaleUp.Step})
	c.Assert(al.Enabled, check.Equals, true)
	c.Assert(al.Actions, check.DeepEquals, []string{"scale_up"})
	scaleName = "scale_down_test"
	al, err = alarm.FindAlarmByName(scaleName)
	c.Assert(err, check.IsNil)
	c.Assert(al.Name, check.Equals, scaleName)
	c.Assert(al.Expression, check.Equals, fmt.Sprintf("%s %s %s", scaleDown.Metric, scaleDown.Operator, scaleDown.Value))
	c.Assert(al.Envs, check.DeepEquals, map[string]string{"step": scaleDown.Step})
	c.Assert(al.Enabled, check.Equals, true)
	c.Assert(al.Actions, check.DeepEquals, []string{"scale_down"})
	alarmName := fmt.Sprintf("enable_scale_down_%s", a.Name)
	al, err = alarm.FindAlarmByName(alarmName)
	c.Assert(err, check.IsNil)
	c.Assert(al.Name, check.Equals, alarmName)
	c.Assert(al.Expression, check.Equals, fmt.Sprintf("units > %d", a.MinUnits))
	c.Assert(al.Envs, check.DeepEquals, map[string]string{"alarm": fmt.Sprintf("scale_down_%s", a.Name)})
	c.Assert(al.Enabled, check.Equals, true)
	c.Assert(al.Actions, check.DeepEquals, []string{"enable_alarm"})
	alarmName = fmt.Sprintf("disable_scale_down_%s", a.Name)
	al, err = alarm.FindAlarmByName(alarmName)
	c.Assert(err, check.IsNil)
	c.Assert(al.Name, check.Equals, alarmName)
	c.Assert(al.Expression, check.Equals, fmt.Sprintf("units <= %d", a.MinUnits))
	c.Assert(al.Envs, check.DeepEquals, map[string]string{"alarm": fmt.Sprintf("scale_down_%s", a.Name)})
	c.Assert(al.Enabled, check.Equals, true)
	c.Assert(al.Actions, check.DeepEquals, []string{"disable_alarm"})
	var as AutoScale
	err = s.conn.Wizard().Find(&a).One(&as)
	c.Assert(err, check.IsNil)
	c.Assert(as.Name, check.Equals, a.Name)
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

func (s *S) TestDisableScaleDown(c *check.C) {
	minUnits := 2
	instanceName := "instanceName"
	err := disableScaleDown(instanceName, minUnits)
	c.Assert(err, check.IsNil)
	alarmName := fmt.Sprintf("disable_scale_down_%s", instanceName)
	al, err := alarm.FindAlarmByName(alarmName)
	c.Assert(err, check.IsNil)
	c.Assert(al.Name, check.Equals, alarmName)
	c.Assert(al.Expression, check.Equals, fmt.Sprintf("units <= %d", minUnits))
	c.Assert(al.Envs, check.DeepEquals, map[string]string{"alarm": fmt.Sprintf("scale_down_%s", instanceName)})
	c.Assert(al.Enabled, check.Equals, true)
	c.Assert(al.Actions, check.DeepEquals, []string{"disable_alarm"})
}

func (s *S) TestAutoScaleUnmarshal(c *check.C) {
	data := []byte(`{"name":"test","minUnits":2,"scaleUp":{},"scaleDown":{}}`)
	a := &AutoScale{}
	err := json.Unmarshal(data, a)
	c.Assert(err, check.IsNil)
}

func (s *S) TestScaleActionUnmarshal(c *check.C) {
	data := []byte(`{"metric":"cpu","operator":">","value":"10","step":"2","wait":200}`)
	sa := &scaleAction{}
	err := json.Unmarshal(data, sa)
	c.Assert(err, check.IsNil)
}

func (s *S) TestFindByName(c *check.C) {
	a := AutoScale{
		Name: "xpto123",
	}
	s.conn.Wizard().Insert(&a)
	a = AutoScale{
		Name: "xpto1234",
	}
	s.conn.Wizard().Insert(&a)
	na, err := FindByName(a.Name)
	c.Assert(err, check.IsNil)
	c.Assert(na, check.DeepEquals, &a)
}
