// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/tsuru/tsuru-autoscale/alarm"
	"gopkg.in/check.v1"
)

func (s *S) TestNewAlarm(c *check.C) {
	body := `{"name":"new","url":"http://tsuru.io","method":"GET"}`
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/alarm", strings.NewReader(body))
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusCreated)
}

func (s *S) TestListAlarms(c *check.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"Name":"instance"}]`))
	}))
	defer ts.Close()
	err := os.Setenv("TSURU_HOST", ts.URL)
	c.Assert(err, check.IsNil)
	err = alarm.NewAlarm(&alarm.Alarm{Name: "myalarm", Instance: "instance"})
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/alarm", nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	c.Assert(recorder.HeaderMap["Content-Type"], check.DeepEquals, []string{"application/json"})
	body := recorder.Body.Bytes()
	var a []alarm.Alarm
	err = json.Unmarshal(body, &a)
	c.Assert(err, check.IsNil)
	c.Assert(a, check.HasLen, 1)
}

func (s *S) TestRemoveAlarmNotFound(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("DELETE", "/alarm/notfound", nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusNotFound)
}

func (s *S) TestRemoveAlarm(c *check.C) {
	a := &alarm.Alarm{Name: "myalarm"}
	err := alarm.NewAlarm(a)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("DELETE", fmt.Sprintf("/alarm/%s", a.Name), nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	_, err = alarm.FindAlarmByName(a.Name)
	c.Assert(err, check.NotNil)
}

func (s *S) TestEnableAlarm(c *check.C) {
	a := &alarm.Alarm{Name: "myalarm", Enabled: false}
	err := alarm.NewAlarm(a)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("PUT", fmt.Sprintf("/alarm/%s/enable", a.Name), nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	a, err = alarm.FindAlarmByName(a.Name)
	c.Assert(err, check.IsNil)
	c.Assert(a.Enabled, check.Equals, true)
}

func (s *S) TestDisableAlarm(c *check.C) {
	a := &alarm.Alarm{Name: "myalarm", Enabled: true}
	err := alarm.NewAlarm(a)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("PUT", fmt.Sprintf("/alarm/%s/disable", a.Name), nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	a, err = alarm.FindAlarmByName(a.Name)
	c.Assert(err, check.IsNil)
	c.Assert(a.Enabled, check.Equals, false)
}

func (s *S) TestGetAlarm(c *check.C) {
	a := &alarm.Alarm{Name: "myalarm"}
	err := alarm.NewAlarm(a)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", fmt.Sprintf("/alarm/%s", a.Name), nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	c.Assert(recorder.HeaderMap["Content-Type"], check.DeepEquals, []string{"application/json"})
	body := recorder.Body.Bytes()
	var got alarm.Alarm
	err = json.Unmarshal(body, &got)
	c.Assert(err, check.IsNil)
	c.Assert(a.Name, check.Equals, got.Name)
}

func (s *S) TestListEvents(c *check.C) {
	a := &alarm.Alarm{Name: "myalarm"}
	err := alarm.NewAlarm(a)
	c.Assert(err, check.IsNil)
	_, err = alarm.NewEvent(a, nil)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/alarm/myalarm/event", nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	c.Assert(recorder.HeaderMap["Content-Type"], check.DeepEquals, []string{"application/json"})
	body := recorder.Body.Bytes()
	var events []alarm.Event
	err = json.Unmarshal(body, &events)
	c.Assert(err, check.IsNil)
	c.Assert(events, check.HasLen, 1)
}

func (s *S) TestListAlarmsByInstance(c *check.C) {
	err := alarm.NewAlarm(&alarm.Alarm{Name: "myalarm", Instance: "instance"})
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/alarm/instance/instance", nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	c.Assert(recorder.HeaderMap["Content-Type"], check.DeepEquals, []string{"application/json"})
	body := recorder.Body.Bytes()
	var a []alarm.Alarm
	err = json.Unmarshal(body, &a)
	c.Assert(err, check.IsNil)
	c.Assert(a, check.HasLen, 1)
}
