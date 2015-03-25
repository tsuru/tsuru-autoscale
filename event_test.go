// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"

	"gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
)

func (s *S) TestLastScaleEvent(c *check.C) {
	config := &Config{Name: "newconfig"}
	event1, err := NewEvent(config, "increase")
	c.Assert(err, check.IsNil)
	event1.StartTime = event1.StartTime.Add(-1 * time.Hour)
	err = event1.update(nil)
	c.Assert(err, check.IsNil)
	event2, err := NewEvent(config, "increase")
	c.Assert(err, check.IsNil)
	event, err := lastScaleEvent(config)
	c.Assert(err, check.IsNil)
	c.Assert(event.ID, check.DeepEquals, event2.ID)
}

func (s *S) TestLastScaleEventNotFound(c *check.C) {
	config := &Config{Name: "not found"}
	_, err := lastScaleEvent(config)
	c.Assert(err, check.Equals, mgo.ErrNotFound)
}

func (s *S) TestEventsByConfigNameWithoutName(c *check.C) {
	config := Config{Name: "config"}
	_, err := NewEvent(&config, "increase")
	c.Assert(err, check.IsNil)
	events, err := eventsByConfigName(nil)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
	c.Assert(events[0].Type, check.Equals, "increase")
	c.Assert(events[0].StartTime, check.Not(check.DeepEquals), time.Time{})
}

func (s *S) TestEventsByConfigName(c *check.C) {
	config := Config{Name: "config"}
	_, err := NewEvent(&config, "increase")
	c.Assert(err, check.IsNil)
	config = Config{Name: "another"}
	_, err = NewEvent(&config, "increase")
	c.Assert(err, check.IsNil)
	events, err := eventsByConfigName(&config)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
	c.Assert(events[0].Type, check.Equals, "increase")
	c.Assert(events[0].StartTime, check.Not(check.DeepEquals), time.Time{})
}
