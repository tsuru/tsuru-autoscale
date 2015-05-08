// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package alarm

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/tsuru/tsuru/db/dbtest"
	"github.com/tsuru/tsuru-autoscale/action"
	"github.com/tsuru/tsuru-autoscale/datasource"
	"github.com/tsuru/tsuru-autoscale/db"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct {
	conn *db.Storage
}

func (s *S) SetUpSuite(c *check.C) {
	var err error
	s.conn, err = db.Conn()
	c.Assert(err, check.IsNil)
}

func (s *S) TearDownTest(c *check.C) {
        dbtest.ClearAllCollections(s.conn.Alarms().Database)
}

var _ = check.Suite(&S{})

func (s *S) TestAlarm(c *check.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"id":"ble"}`))
	}))
	defer ts.Close()
	instance := datasource.DataSource{
		Name:   "ds",
		URL:    ts.URL,
		Method: "GET",
	}
	alarm := Alarm{
		Name:       "name",
		Expression: `data.id == "ble"`,
		DataSource: instance,
	}
	err := NewAlarm(&alarm)
	c.Assert(err, check.IsNil)
	err = scaleIfNeeded(&alarm)
	c.Assert(err, check.IsNil)
	var events []Event
	err = s.conn.Events().Find(nil).All(&events)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
}

func (s *S) TestRunAutoScaleOnce(c *check.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"id":"ble"}`))
	}))
	defer ts.Close()
	instance := datasource.DataSource{
		Name:   "ds",
		URL:    ts.URL,
		Method: "GET",
	}
	alarm := Alarm{
		Name:       "name",
		Expression: `data.id == "ble"`,
		DataSource: instance,
	}
	err := NewAlarm(&alarm)
	c.Assert(err, check.IsNil)
	runAutoScaleOnce()
	var events []Event
	err = s.conn.Events().Find(nil).All(&events)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
	c.Assert(events[0].StartTime, check.Not(check.DeepEquals), time.Time{})
	c.Assert(events[0].EndTime, check.Not(check.DeepEquals), time.Time{})
	c.Assert(events[0].Error, check.Equals, "")
	c.Assert(events[0].Successful, check.Equals, true)
}

func (s *S) TestAutoScaleEnable(c *check.C) {
	alarm := Alarm{Name: "alarm"}
	err := AutoScaleEnable(&alarm)
	c.Assert(err, check.IsNil)
	c.Assert(alarm.Enabled, check.Equals, true)
}

func (s *S) TestAutoScaleDisable(c *check.C) {
	alarm := Alarm{Name: "alarm", Enabled: true}
	err := AutoScaleDisable(&alarm)
	c.Assert(err, check.IsNil)
	c.Assert(alarm.Enabled, check.Equals, false)
}

func (s *S) TestAlarmWaitEventStillRunning(c *check.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"id":"ble"}`))
	}))
	defer ts.Close()
	instance := datasource.DataSource{
		Name:   "ds",
		URL:    ts.URL,
		Method: "GET",
	}
	a := action.Action{
		Name:   "name",
		Method: "GET",
		URL:    "http:/tsuru.io",
	}
	alarm := &Alarm{
		Name:       "rush",
		Actions:    []action.Action{a},
		Enabled:    true,
		DataSource: instance,
	}
	event, err := NewEvent(alarm)
	c.Assert(err, check.IsNil)
	err = scaleIfNeeded(alarm)
	c.Assert(err, check.IsNil)
	events, err := eventsByAlarmName(alarm)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
	c.Assert(events[0].ID, check.DeepEquals, event.ID)
}

func (s *S) TestAlarmWaitTime(c *check.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"id":"ble"}`))
	}))
	defer ts.Close()
	instance := datasource.DataSource{
		Name:   "ds",
		URL:    ts.URL,
		Method: "GET",
	}
	a := action.Action{
		Name:   "name",
		URL:    "http://tsuru.io",
		Method: "GET",
	}
	alarm := &Alarm{
		Name:       "rush",
		Actions:    []action.Action{a},
		Enabled:    true,
		DataSource: instance,
	}
	event, err := NewEvent(alarm)
	c.Assert(err, check.IsNil)
	err = event.update(nil)
	c.Assert(err, check.IsNil)
	err = scaleIfNeeded(alarm)
	c.Assert(err, check.IsNil)
	events, err := eventsByAlarmName(alarm)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
	c.Assert(events[0].ID, check.DeepEquals, event.ID)
}

func (s *S) TestAlarmCheck(c *check.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"id":"ble"}`))
	}))
	defer ts.Close()
	instance := datasource.DataSource{
		Name:   "ds",
		URL:    ts.URL,
		Method: "GET",
	}
	alarm := &Alarm{
		Name:       "rush",
		Enabled:    true,
		Expression: `data.id == "ble"`,
		DataSource: instance,
	}
	ok, err := alarm.Check()
	c.Assert(err, check.IsNil)
	c.Assert(ok, check.Equals, true)
	alarm = &Alarm{
		Name:       "rush",
		Enabled:    true,
		Expression: `data.id != "ble"`,
		DataSource: instance,
	}
	ok, err = alarm.Check()
	c.Assert(err, check.IsNil)
	c.Assert(ok, check.Equals, false)
}

func (s *S) TestListAlarms(c *check.C) {
	a := Alarm{
		Name: "xpto",
	}
	s.conn.Alarms().Insert(&a)
	a = Alarm{
		Name: "xpto2",
	}
	s.conn.Alarms().Insert(&a)
	all, err := ListAlarms()
	c.Assert(err, check.IsNil)
	c.Assert(all, check.HasLen, 2)
}
