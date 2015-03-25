// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

func (s *S) TestEvent(c *check.C) {
	strg, err := Conn()
	c.Assert(err, check.IsNil)
	event := strg.Event()
	eventc := strg.Collection("event")
	c.Assert(event, check.DeepEquals, eventc)
}

func (s *S) TestConfig(c *check.C) {
	strg, err := Conn()
	c.Assert(err, check.IsNil)
	config := strg.Config()
	configc := strg.Collection("config")
	c.Assert(config, check.DeepEquals, configc)
}
