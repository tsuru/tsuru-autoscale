// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"net/http"
	"net/http/httptest"

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
