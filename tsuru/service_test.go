// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsuru

import (
	"testing"

	"github.com/tsuru/tsuru-autoscale/db"
	"github.com/tsuru/tsuru/db/dbtest"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct {
	conn *db.Storage
}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {
	var err error
	s.conn, err = db.Conn()
	c.Assert(err, check.IsNil)
}
func (s *S) TearDownTest(c *check.C) {
	dbtest.ClearAllCollections(s.conn.Instances().Database)
}

func (s *S) TestNewInstance(c *check.C) {
	i := &Instance{
		Name: "name",
	}
	err := NewInstance(i)
	c.Assert(err, check.IsNil)
}

func (s *S) TestGetInstanceByName(c *check.C) {
	i := &Instance{
		Name: "name",
	}
	err := NewInstance(i)
	c.Assert(err, check.IsNil)
	n, err := GetInstanceByName(i.Name)
	c.Assert(err, check.IsNil)
	c.Assert(n.Name, check.Equals, i.Name)
}

func (s *S) TestAddApp(c *check.C) {
	i := &Instance{
		Name: "name",
	}
	err := NewInstance(i)
	c.Assert(err, check.IsNil)
	i, err = GetInstanceByName(i.Name)
	c.Assert(err, check.IsNil)
	err = i.AddApp("app.domain.com")
	c.Assert(err, check.IsNil)
	err = i.AddApp("app.domain.com")
	c.Assert(err, check.NotNil)
	i, err = GetInstanceByName(i.Name)
	c.Assert(err, check.IsNil)
	c.Assert(i.Apps, check.DeepEquals, []string{"app"})
}

func (s *S) TestRemoveInstance(c *check.C) {
	i := &Instance{
		Name: "name",
	}
	err := NewInstance(i)
	c.Assert(err, check.IsNil)
	n, err := GetInstanceByName(i.Name)
	c.Assert(err, check.IsNil)
	c.Assert(n.Name, check.Equals, i.Name)
	err = RemoveInstance(n)
	c.Assert(err, check.IsNil)
	n, err = GetInstanceByName(i.Name)
	c.Assert(err, check.NotNil)
	c.Assert(n, check.IsNil)
}
