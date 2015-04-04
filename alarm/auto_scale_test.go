// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package alarm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/tsuru/config"
	"github.com/tsuru/tsuru-autoscale/action"
	"github.com/tsuru/tsuru-autoscale/db"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct {
	conn *db.Storage
}

func (s *S) SetUpSuite(c *check.C) {
	err := config.ReadConfigFile("testdata/config.yaml")
	c.Assert(err, check.IsNil)
	s.conn, err = db.Conn()
	c.Assert(err, check.IsNil)
}

func (s *S) TearDownTest(c *check.C) {
	s.conn.Events().RemoveAll(nil)
}

var _ = check.Suite(&S{})

type metricHandler struct {
	cpuMax string
}

func (h *metricHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	content := fmt.Sprintf(`[{"target": "sometarget", "datapoints": [[2.2, 1415129040], [2.2, 1415129050], [2.2, 1415129060], [2.2, 1415129070], [%s, 1415129080]]}]`, h.cpuMax)
	w.Write([]byte(content))
}

func (s *S) TestAlarm(c *check.C) {
	h := metricHandler{cpuMax: "50.2"}
	ts := httptest.NewServer(&h)
	defer ts.Close()
	alarm := &Alarm{
		Action:  action.Action{Units: 1, Expression: "{cpu} > 80"},
		Enabled: true,
	}
	err := scaleIfNeeded(alarm)
	c.Assert(err, check.IsNil)
	var events []Event
	err = s.conn.Events().Find(nil).All(&events)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 0)
}

type autoscaleHandler struct {
	matches map[string]string
}

func (h *autoscaleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var cpu string
	for key, value := range h.matches {
		if strings.Contains(r.URL.String(), key) {
			cpu = value
		}
	}
	content := fmt.Sprintf(`[{"target": "sometarget", "datapoints": [[2.2, 1415129040], [2.2, 1415129050], [2.2, 1415129060], [2.2, 1415129070], [%s, 1415129080]]}]`, cpu)
	w.Write([]byte(content))
}

func (s *S) TestRunAutoScaleOnce(c *check.C) {
	h := autoscaleHandler{
		matches: map[string]string{
			"myApp":      "90.2",
			"anotherApp": "9.2",
		},
	}
	ts := httptest.NewServer(&h)
	defer ts.Close()
	up := &Alarm{
		Action:  action.Action{Units: 1, Expression: "{cpu_max} > 80"},
		Enabled: true,
	}
	dh := metricHandler{cpuMax: "9.2"}
	dts := httptest.NewServer(&dh)
	defer dts.Close()
	down := &Alarm{
		Action:  action.Action{Units: 1, Expression: "{cpu_max} > 80"},
		Enabled: true,
	}
	runAutoScaleOnce()
	var events []Event
	err := s.conn.Events().Find(nil).All(&events)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 2)
	c.Assert(events[0].Type, check.Equals, "increase")
	c.Assert(events[0].StartTime, check.Not(check.DeepEquals), time.Time{})
	c.Assert(events[0].EndTime, check.Not(check.DeepEquals), time.Time{})
	c.Assert(events[0].Error, check.Equals, "")
	c.Assert(events[0].Successful, check.Equals, true)
	c.Assert(events[0].Alarm, check.DeepEquals, up)
	c.Assert(events[1].Type, check.Equals, "decrease")
	c.Assert(events[1].StartTime, check.Not(check.DeepEquals), time.Time{})
	c.Assert(events[1].EndTime, check.Not(check.DeepEquals), time.Time{})
	c.Assert(events[1].Error, check.Equals, "")
	c.Assert(events[1].Successful, check.Equals, true)
	c.Assert(events[1].Alarm, check.DeepEquals, down)
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
	h := metricHandler{cpuMax: "10.2"}
	ts := httptest.NewServer(&h)
	defer ts.Close()
	alarm := &Alarm{
		Name:    "rush",
		Action:  action.Action{Units: 5, Expression: "{cpu_max} > 80", Wait: 30e9},
		Enabled: true,
	}
	event, err := NewEvent(alarm, "decrease")
	c.Assert(err, check.IsNil)
	err = scaleIfNeeded(alarm)
	c.Assert(err, check.IsNil)
	events, err := eventsByAlarmName(alarm)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
	c.Assert(events[0].ID, check.DeepEquals, event.ID)
}

func (s *S) TestAlarmWaitTime(c *check.C) {
	h := metricHandler{cpuMax: "10.2"}
	ts := httptest.NewServer(&h)
	defer ts.Close()
	alarm := &Alarm{
		Name:    "rush",
		Action:  action.Action{Units: 5, Expression: "{cpu_max} > 80", Wait: 1 * time.Hour},
		Enabled: true,
	}
	event, err := NewEvent(alarm, "increase")
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

func (s *S) TestAlarmMarshalJSON(c *check.C) {
	conf := &Alarm{
		Action:  action.Action{Units: 1, Expression: "{cpu} > 80"},
		Enabled: true,
	}
	expected := map[string]interface{}{
		"name":       "",
		"expression": "",
		"wait":       float64(0),
		"action": map[string]interface{}{
			"wait":       float64(0),
			"expression": "{cpu} > 80",
			"units":      float64(1),
		},
		"enabled": true,
	}
	data, err := json.Marshal(conf)
	c.Assert(err, check.IsNil)
	result := make(map[string]interface{})
	err = json.Unmarshal(data, &result)
	c.Assert(err, check.IsNil)
	c.Assert(result, check.DeepEquals, expected)
}
