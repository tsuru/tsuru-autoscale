// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/ajg/form"
	"github.com/tsuru/tsuru-autoscale/alarm"
	"gopkg.in/check.v1"
)

func (s *S) TestAlarmEnable(c *check.C) {
	a := &alarm.Alarm{Name: "myalarm", Enabled: false}
	err := alarm.NewAlarm(a)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/alarm/myalarm/enable", nil)
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusFound)
	a, err = alarm.FindAlarmByName("myalarm")
	c.Assert(err, check.IsNil)
	c.Assert(a.Enabled, check.Equals, true)
}

func (s *S) TestAlarmDisable(c *check.C) {
	a := &alarm.Alarm{Name: "myalarm", Enabled: true}
	err := alarm.NewAlarm(a)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/alarm/myalarm/disable", nil)
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusFound)
	a, err = alarm.FindAlarmByName("myalarm")
	c.Assert(err, check.IsNil)
	c.Assert(a.Enabled, check.Equals, false)
}

func (s *S) TestAlarmRemove(c *check.C) {
	a := &alarm.Alarm{Name: "myalarm"}
	err := alarm.NewAlarm(a)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/alarm/myalarm/delete", nil)
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusFound)
	_, err = alarm.FindAlarmByName("myalarm")
	c.Assert(err, check.NotNil)
}

func (s *S) TestAlarmAdd(c *check.C) {
	v := url.Values{
		"key":         []string{"", "f", "x"},
		"value":       []string{"", "f", "x"},
		"name":        []string{"new"},
		"enabled":     []string{"true"},
		"datasources": []string{"cpu", "memory"},
		"actions":     []string{"up", "down"},
	}
	body := strings.NewReader(v.Encode())
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/alarm/add", body)
	c.Assert(err, check.IsNil)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	server(recorder, request)
	c.Assert(recorder.Body.String(), check.Equals, "")
	c.Assert(recorder.Code, check.Equals, http.StatusFound)
	r, err := alarm.FindAlarmByName("new")
	c.Assert(err, check.IsNil)
	c.Assert(r.Name, check.Equals, "new")
	c.Assert(r.DataSources, check.DeepEquals, []string{"cpu", "memory"})
	c.Assert(r.Actions, check.DeepEquals, []string{"up", "down"})
	c.Assert(r.Envs, check.DeepEquals, map[string]string{"x": "x", "f": "f"})
}

func (s *S) TestAlarmEdit(c *check.C) {
	a := &alarm.Alarm{Name: "myalarm", Enabled: true, Instance: "myalarm-instance"}
	err := alarm.NewAlarm(a)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	v := url.Values{
		"key":         []string{"", "f", "x"},
		"value":       []string{"", "f", "x"},
		"name":        []string{"myalarm"},
		"enabled":     []string{"false"},
		"datasources": []string{"cpu", "memory"},
		"actions":     []string{"up", "down"},
	}
	body := strings.NewReader(v.Encode())
	request, err := http.NewRequest("POST", "/alarm/myalarm/edit", body)
	c.Assert(err, check.IsNil)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusFound)
	r, err := alarm.FindAlarmByName(a.Name)
	c.Assert(err, check.IsNil)
	c.Assert(r.Enabled, check.Equals, false)
	c.Assert(r.DataSources, check.DeepEquals, []string{"cpu", "memory"})
	c.Assert(r.Actions, check.DeepEquals, []string{"up", "down"})
	c.Assert(r.Envs, check.DeepEquals, map[string]string{"x": "x", "f": "f"})
	c.Assert(r.Instance, check.Equals, "myalarm-instance")
}

func (s *S) TestAlarmEditEmptyBody(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/alarm/myalarm/edit", nil)
	c.Assert(err, check.IsNil)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusInternalServerError)
}

func (s *S) TestAlarmEditNotFound(c *check.C) {
	recorder := httptest.NewRecorder()
	a := &alarm.Alarm{Name: "myalarm", Enabled: false}
	v, err := form.EncodeToValues(&a)
	c.Assert(err, check.IsNil)
	body := strings.NewReader(v.Encode())
	request, err := http.NewRequest("POST", "/alarm/myalarm/edit", body)
	c.Assert(err, check.IsNil)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusNotFound)
}
