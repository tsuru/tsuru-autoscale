// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/tsuru/tsuru-autoscale/alarm"
	"github.com/tsuru/tsuru-autoscale/wizard"
	"gopkg.in/check.v1"
)

func (s *S) TestNewAutoScale(c *check.C) {
	body := `{"name":"test","minUnits":2,"scaleUp":{},"scaleDown":{}}`
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/wizard", strings.NewReader(body))
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusCreated)
}

func (s *S) TestWizardByName(c *check.C) {
	autoScale := &wizard.AutoScale{
		Name: "instance",
	}
	err := wizard.New(autoScale)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/wizard/instance", nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	c.Assert(recorder.HeaderMap["Content-Type"], check.DeepEquals, []string{"application/json"})
	body := recorder.Body.Bytes()
	var instance wizard.AutoScale
	err = json.Unmarshal(body, &instance)
	c.Assert(err, check.IsNil)
	c.Assert(instance.Name, check.Equals, "instance")
}

func (s *S) TestRemoveWizardNotFound(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("DELETE", "/wizard/notfound", nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusNotFound)
}

func (s *S) TestRemoveWizard(c *check.C) {
	scaleUp := wizard.ScaleAction{
		Metric:   "cpu",
		Operator: ">",
		Step:     "1",
		Value:    "10",
		Wait:     50,
	}
	scaleDown := wizard.ScaleAction{
		Metric:   "cpu",
		Operator: "<",
		Step:     "1",
		Value:    "2",
		Wait:     50,
	}
	autoScale := &wizard.AutoScale{
		Name:      "instance",
		ScaleUp:   scaleUp,
		ScaleDown: scaleDown,
		Process:   "web",
	}
	err := wizard.New(autoScale)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("DELETE", fmt.Sprintf("/wizard/%s", autoScale.Name), nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	_, err = wizard.FindByName(autoScale.Name)
	c.Assert(err, check.NotNil)
}

func (s *S) TestEventsByWizardName(c *check.C) {
	al := alarm.Alarm{
		Name:     "enable_scale_down_xpto1234",
		Instance: "xpto1234",
		Actions:  []string{"scale_down"},
	}
	_, err := alarm.NewEvent(&al, nil)
	c.Assert(err, check.IsNil)
	a := wizard.AutoScale{
		Name: "xpto1234",
	}
	err = wizard.New(&a)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/wizard/xpto1234/events", nil)
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

func (s *S) TestEnableWizardNotFound(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/wizard/notfound/enable", nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusNotFound)
}

func (s *S) TestEnableWizard(c *check.C) {
	scaleUp := wizard.ScaleAction{
		Metric:   "cpu",
		Operator: ">",
		Step:     "1",
		Value:    "10",
		Wait:     50,
	}
	scaleDown := wizard.ScaleAction{
		Metric:   "cpu",
		Operator: "<",
		Step:     "1",
		Value:    "2",
		Wait:     50,
	}
	autoScale := &wizard.AutoScale{
		Name:      "instance",
		ScaleUp:   scaleUp,
		ScaleDown: scaleDown,
		Process:   "web",
	}
	err := wizard.New(autoScale)
	c.Assert(err, check.IsNil)
	err = autoScale.Disable()
	c.Assert(err, check.IsNil)
	c.Assert(autoScale.Enabled(), check.Equals, false)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/wizard/%s/enable", autoScale.Name), nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	a, err := wizard.FindByName(autoScale.Name)
	c.Assert(err, check.IsNil)
	c.Assert(a.Enabled(), check.Equals, true)
}
