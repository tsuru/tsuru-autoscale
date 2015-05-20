// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package alarm

import (
	"time"

	"gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
)

func (s *S) TestLastScaleEvent(c *check.C) {
	alarm := &Alarm{Name: "newconfig"}
	event1, err := NewEvent(alarm, nil)
	c.Assert(err, check.IsNil)
	event1.StartTime = event1.StartTime.Add(-1 * time.Hour)
	err = event1.update(nil)
	c.Assert(err, check.IsNil)
	event2, err := NewEvent(alarm, nil)
	c.Assert(err, check.IsNil)
	event, err := lastScaleEvent(alarm)
	c.Assert(err, check.IsNil)
	c.Assert(event.ID, check.DeepEquals, event2.ID)
}

func (s *S) TestLastScaleEventNotFound(c *check.C) {
	alarm := &Alarm{Name: "not found"}
	_, err := lastScaleEvent(alarm)
	c.Assert(err, check.Equals, mgo.ErrNotFound)
}

func (s *S) TestEventsByAlarmNameWithoutName(c *check.C) {
	alarm := Alarm{Name: "config"}
	_, err := NewEvent(&alarm, nil)
	c.Assert(err, check.IsNil)
	events, err := eventsByAlarmName(nil)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
	c.Assert(events[0].StartTime, check.Not(check.DeepEquals), time.Time{})
}

func (s *S) TestEventsByAlarmName(c *check.C) {
	alarm := Alarm{Name: "config"}
	_, err := NewEvent(&alarm, nil)
	c.Assert(err, check.IsNil)
	alarm = Alarm{Name: "another"}
	_, err = NewEvent(&alarm, nil)
	c.Assert(err, check.IsNil)
	events, err := eventsByAlarmName(&alarm)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
	c.Assert(events[0].StartTime, check.Not(check.DeepEquals), time.Time{})
}
