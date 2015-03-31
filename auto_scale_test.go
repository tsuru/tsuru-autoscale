// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

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

func (s *S) TestAutoScale(c *check.C) {
	h := metricHandler{cpuMax: "50.2"}
	ts := httptest.NewServer(&h)
	defer ts.Close()
	config := &Config{
		Increase: action.Action{Units: 1, Expression: "{cpu} > 80"},
		Decrease: action.Action{Units: 1, Expression: "{cpu} < 20"},
		Enabled:  true,
	}
	err := scaleIfNeeded(config)
	c.Assert(err, check.IsNil)
	var events []Event
	err = s.conn.Events().Find(nil).All(&events)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 0)
}

func (s *S) TestAutoScaleUp(c *check.C) {
	h := metricHandler{cpuMax: "90.2"}
	ts := httptest.NewServer(&h)
	defer ts.Close()
	config := &Config{
		Increase: action.Action{Units: 1, Expression: "{cpu_max} > 80"},
		Enabled:  true,
		MaxUnits: uint(10),
	}
	err := scaleIfNeeded(config)
	c.Assert(err, check.IsNil)
	var events []Event
	err = s.conn.Events().Find(nil).All(&events)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
	c.Assert(events[0].Type, check.Equals, "increase")
	c.Assert(events[0].StartTime, check.Not(check.DeepEquals), time.Time{})
	c.Assert(events[0].EndTime, check.Not(check.DeepEquals), time.Time{})
	c.Assert(events[0].Error, check.Equals, "")
	c.Assert(events[0].Successful, check.Equals, true)
	c.Assert(events[0].Config, check.DeepEquals, config)
}

func (s *S) TestAutoScaleDown(c *check.C) {
	h := metricHandler{cpuMax: "10.2"}
	ts := httptest.NewServer(&h)
	defer ts.Close()
	config := &Config{
		Increase: action.Action{Units: 1, Expression: "{cpu_max} > 80"},
		Decrease: action.Action{Units: 1, Expression: "{cpu_max} < 20"},
		Enabled:  true,
	}
	err := scaleIfNeeded(config)
	c.Assert(err, check.IsNil)
	var events []Event
	err = s.conn.Events().Find(nil).All(&events)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
	c.Assert(events[0].Type, check.Equals, "decrease")
	c.Assert(events[0].StartTime, check.Not(check.DeepEquals), time.Time{})
	c.Assert(events[0].EndTime, check.Not(check.DeepEquals), time.Time{})
	c.Assert(events[0].Error, check.Equals, "")
	c.Assert(events[0].Successful, check.Equals, true)
	c.Assert(events[0].Config, check.DeepEquals, config)
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
	up := &Config{
		Increase: action.Action{Units: 1, Expression: "{cpu_max} > 80"},
		Enabled:  true,
		MaxUnits: uint(10),
	}
	dh := metricHandler{cpuMax: "9.2"}
	dts := httptest.NewServer(&dh)
	defer dts.Close()
	down := &Config{
		Increase: action.Action{Units: 1, Expression: "{cpu_max} > 80"},
		Decrease: action.Action{Units: 1, Expression: "{cpu_max} < 20"},
		Enabled:  true,
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
	c.Assert(events[0].Config, check.DeepEquals, up)
	c.Assert(events[1].Type, check.Equals, "decrease")
	c.Assert(events[1].StartTime, check.Not(check.DeepEquals), time.Time{})
	c.Assert(events[1].EndTime, check.Not(check.DeepEquals), time.Time{})
	c.Assert(events[1].Error, check.Equals, "")
	c.Assert(events[1].Successful, check.Equals, true)
	c.Assert(events[1].Config, check.DeepEquals, down)
}

func (s *S) TestAutoScaleEnable(c *check.C) {
	config := Config{Name: "config"}
	err := AutoScaleEnable(&config)
	c.Assert(err, check.IsNil)
	c.Assert(config.Enabled, check.Equals, true)
}

func (s *S) TestAutoScaleDisable(c *check.C) {
	config := Config{Name: "config", Enabled: true}
	err := AutoScaleDisable(&config)
	c.Assert(err, check.IsNil)
	c.Assert(config.Enabled, check.Equals, false)
}

func (s *S) TestAutoScaleUpWaitEventStillRunning(c *check.C) {
	h := metricHandler{cpuMax: "90.2"}
	ts := httptest.NewServer(&h)
	defer ts.Close()
	config := &Config{
		Increase: action.Action{Units: 5, Expression: "{cpu_max} > 80", Wait: 30e9},
		Enabled:  true,
		MaxUnits: 4,
	}
	event, err := NewEvent(config, "increase")
	c.Assert(err, check.IsNil)
	err = scaleIfNeeded(config)
	c.Assert(err, check.IsNil)
	events, err := eventsByConfigName(config)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
	c.Assert(events[0].ID, check.DeepEquals, event.ID)
}

func (s *S) TestAutoScaleUpWaitTime(c *check.C) {
	h := metricHandler{cpuMax: "90.2"}
	ts := httptest.NewServer(&h)
	defer ts.Close()
	config := &Config{
		Increase: action.Action{Units: 5, Expression: "{cpu_max} > 80", Wait: 1 * time.Hour},
		Enabled:  true,
		MaxUnits: 4,
	}
	event, err := NewEvent(config, "increase")
	c.Assert(err, check.IsNil)
	err = event.update(nil)
	c.Assert(err, check.IsNil)
	err = scaleIfNeeded(config)
	c.Assert(err, check.IsNil)
	events, err := eventsByConfigName(config)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
	c.Assert(events[0].ID, check.DeepEquals, event.ID)
}

func (s *S) TestAutoScaleMaxUnits(c *check.C) {
	h := metricHandler{cpuMax: "90.2"}
	ts := httptest.NewServer(&h)
	defer ts.Close()
	config := &Config{
		Increase: action.Action{Units: 5, Expression: "{cpu_max} > 80"},
		Enabled:  true,
		MaxUnits: 4,
	}
	err := scaleIfNeeded(config)
	c.Assert(err, check.IsNil)
	var events []Event
	c.Assert(events, check.HasLen, 1)
	c.Assert(events[0].Type, check.Equals, "increase")
	c.Assert(events[0].StartTime, check.Not(check.DeepEquals), time.Time{})
	c.Assert(events[0].EndTime, check.Not(check.DeepEquals), time.Time{})
	c.Assert(events[0].Error, check.Equals, "")
	c.Assert(events[0].Successful, check.Equals, true)
	c.Assert(events[0].Config, check.DeepEquals, config)
}

func (s *S) TestAutoScaleDownWaitEventStillRunning(c *check.C) {
	h := metricHandler{cpuMax: "10.2"}
	ts := httptest.NewServer(&h)
	defer ts.Close()
	config := &Config{
		Name:     "rush",
		Increase: action.Action{Units: 5, Expression: "{cpu_max} > 80", Wait: 30e9},
		Decrease: action.Action{Units: 3, Expression: "{cpu_max} < 20", Wait: 30e9},
		Enabled:  true,
		MaxUnits: 4,
	}
	event, err := NewEvent(config, "decrease")
	c.Assert(err, check.IsNil)
	err = scaleIfNeeded(config)
	c.Assert(err, check.IsNil)
	events, err := eventsByConfigName(config)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
	c.Assert(events[0].ID, check.DeepEquals, event.ID)
}

func (s *S) TestAutoScaleDownWaitTime(c *check.C) {
	h := metricHandler{cpuMax: "10.2"}
	ts := httptest.NewServer(&h)
	defer ts.Close()
	config := &Config{
		Name:     "rush",
		Increase: action.Action{Units: 5, Expression: "{cpu_max} > 80", Wait: 1 * time.Hour},
		Decrease: action.Action{Units: 3, Expression: "{cpu_max} < 20", Wait: 3 * time.Hour},
		Enabled:  true,
		MaxUnits: 4,
	}
	event, err := NewEvent(config, "increase")
	c.Assert(err, check.IsNil)
	err = event.update(nil)
	c.Assert(err, check.IsNil)
	err = scaleIfNeeded(config)
	c.Assert(err, check.IsNil)
	events, err := eventsByConfigName(config)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
	c.Assert(events[0].ID, check.DeepEquals, event.ID)
}

func (s *S) TestAutoScaleMinUnits(c *check.C) {
	h := metricHandler{cpuMax: "10.2"}
	ts := httptest.NewServer(&h)
	defer ts.Close()
	config := &Config{
		Increase: action.Action{Units: 1, Expression: "{cpu_max} > 80"},
		Decrease: action.Action{Units: 3, Expression: "{cpu_max} < 20"},
		Enabled:  true,
		MinUnits: uint(3),
	}
	err := scaleIfNeeded(config)
	c.Assert(err, check.IsNil)
	var events []Event
	err = s.conn.Events().Find(nil).All(&events)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
	c.Assert(events[0].Type, check.Equals, "decrease")
	c.Assert(events[0].StartTime, check.Not(check.DeepEquals), time.Time{})
	c.Assert(events[0].EndTime, check.Not(check.DeepEquals), time.Time{})
	c.Assert(events[0].Error, check.Equals, "")
	c.Assert(events[0].Successful, check.Equals, true)
	c.Assert(events[0].Config, check.DeepEquals, config)
}

func (s *S) TestConfigMarshalJSON(c *check.C) {
	conf := &Config{
		Increase: action.Action{Units: 1, Expression: "{cpu} > 80"},
		Decrease: action.Action{Units: 1, Expression: "{cpu} < 20"},
		Enabled:  true,
		MaxUnits: 10,
		MinUnits: 2,
	}
	expected := map[string]interface{}{
		"name": "",
		"increase": map[string]interface{}{
			"wait":       float64(0),
			"expression": "{cpu} > 80",
			"units":      float64(1),
		},
		"decrease": map[string]interface{}{
			"wait":       float64(0),
			"expression": "{cpu} < 20",
			"units":      float64(1),
		},
		"minUnits": float64(2),
		"maxUnits": float64(10),
		"enabled":  true,
	}
	data, err := json.Marshal(conf)
	c.Assert(err, check.IsNil)
	result := make(map[string]interface{})
	err = json.Unmarshal(data, &result)
	c.Assert(err, check.IsNil)
	c.Assert(result, check.DeepEquals, expected)
}

func (s *S) TestAutoScaleDownMin(c *check.C) {
	h := metricHandler{cpuMax: "10.2"}
	ts := httptest.NewServer(&h)
	defer ts.Close()
	config := &Config{
		Increase: action.Action{Units: 1, Expression: "{cpu_max} > 80"},
		Decrease: action.Action{Units: 1, Expression: "{cpu_max} < 20"},
		Enabled:  true,
		MinUnits: 1,
	}
	err := scaleIfNeeded(config)
	c.Assert(err, check.IsNil)
	var events []Event
	err = s.conn.Events().Find(nil).All(&events)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 0)
}

func (s *S) TestAutoScaleUpMax(c *check.C) {
	h := metricHandler{cpuMax: "90.2"}
	ts := httptest.NewServer(&h)
	defer ts.Close()
	config := &Config{
		Increase: action.Action{Units: 1, Expression: "{cpu_max} > 80"},
		Enabled:  true,
		MaxUnits: uint(2),
	}
	err := scaleIfNeeded(config)
	c.Assert(err, check.IsNil)
	var events []Event
	err = s.conn.Events().Find(nil).All(&events)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 0)
}
