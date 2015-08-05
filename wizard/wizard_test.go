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
		waitTime: 50,
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
