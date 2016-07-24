// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package alarm

import (
	"time"

	"gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	alarm := Alarm{Name: "all"}
	_, err := NewEvent(&alarm, nil)
	c.Assert(err, check.IsNil)
	events, err := EventsByAlarmName("")
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
	c.Assert(events[0].StartTime, check.Not(check.DeepEquals), time.Time{})
}

func (s *S) TestEventsByAlarmOrderByStartTime(c *check.C) {
	alarm := Alarm{Name: "orderedevents"}
	_, err := NewEvent(&alarm, nil)
	c.Assert(err, check.IsNil)
	time.Sleep(1 * time.Second)
	_, err = NewEvent(&alarm, nil)
	c.Assert(err, check.IsNil)
	events, err := EventsByAlarmName("")
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 2)
	c.Assert(events[0].StartTime.After(events[1].StartTime), check.Equals, true)
}

func (s *S) TestEventsByAlarmName(c *check.C) {
	alarm := Alarm{Name: "config"}
	_, err := NewEvent(&alarm, nil)
	c.Assert(err, check.IsNil)
	alarm = Alarm{Name: "another"}
	_, err = NewEvent(&alarm, nil)
	c.Assert(err, check.IsNil)
	events, err := EventsByAlarmName(alarm.Name)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
	c.Assert(events[0].StartTime, check.Not(check.DeepEquals), time.Time{})
}

func (s *S) TestFindEventsBy(c *check.C) {
	alarm := Alarm{Name: "config"}
	_, err := NewEvent(&alarm, nil)
	c.Assert(err, check.IsNil)
	alarm = Alarm{Name: "another"}
	_, err = NewEvent(&alarm, nil)
	c.Assert(err, check.IsNil)
	events, err := FindEventsBy(bson.M{"alarm.name": alarm.Name}, 1000)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
	events, err = FindEventsBy(nil, 1000)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 2)
}
