// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/tsuru/tsuru-autoscale/tsuru"
	"github.com/tsuru/tsuru-autoscale/wizard"
	"gopkg.in/check.v1"
)

func (s *S) TestServiceAdd(c *check.C) {
	recorder := httptest.NewRecorder()
	body := `name=myscale2&team=admin&user=admin%40example.com`
	request, err := http.NewRequest("POST", "/resources", strings.NewReader(body))
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusCreated)
}

func (s *S) TestServiceBindUnit(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/resources/name/bind", nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusCreated)
}

func (s *S) TestServiceBindAppByHost(c *check.C) {
	service := &tsuru.Instance{
		Name: "name",
	}
	err := tsuru.NewInstance(service)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	body := `app-host=tsuru-dashboard.192.168.50.4.nip.io`
	request, err := http.NewRequest("POST", "/resources/name/bind-app", strings.NewReader(body))
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusCreated)
	var i interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &i)
	c.Assert(err, check.IsNil)
}

func (s *S) TestServiceBindApp(c *check.C) {
	service := &tsuru.Instance{
		Name: "name",
	}
	err := tsuru.NewInstance(service)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	body := `app-name=tsuru-dashboard`
	request, err := http.NewRequest("POST", "/resources/name/bind-app", strings.NewReader(body))
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusCreated)
	var i interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &i)
	c.Assert(err, check.IsNil)
}

func (s *S) TestServiceUnbindUnit(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("DELETE", "/resources/name/bind", nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
}

func (s *S) TestServiceUnbindAppByHost(c *check.C) {
	service := &tsuru.Instance{
		Name: "name",
	}
	err := tsuru.NewInstance(service)
	c.Assert(err, check.IsNil)
	instance, err := tsuru.GetInstanceByName("name")
	c.Assert(err, check.IsNil)
	err = instance.AddApp("", "tsuru-dashboard.192.168.50.4.nip.io")
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	body := `app-host=tsuru-dashboard.192.168.50.4.nip.io`
	request, err := http.NewRequest("DELETE", "/resources/name/bind-app", strings.NewReader(body))
	c.Assert(err, check.IsNil)
	request.Header.Add("Authorization", "token")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	instance, err = tsuru.GetInstanceByName("name")
	c.Assert(err, check.IsNil)
	c.Assert(instance.Apps, check.HasLen, 0)
}

func (s *S) TestServiceUnbindApp(c *check.C) {
	service := &tsuru.Instance{
		Name: "name",
	}
	err := tsuru.NewInstance(service)
	c.Assert(err, check.IsNil)
	instance, err := tsuru.GetInstanceByName("name")
	c.Assert(err, check.IsNil)
	err = instance.AddApp("", "tsuru-dashboard.192.168.50.4.nip.io")
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	body := `app-name=tsuru-dashboard`
	request, err := http.NewRequest("DELETE", "/resources/name/bind-app", strings.NewReader(body))
	c.Assert(err, check.IsNil)
	request.Header.Add("Authorization", "token")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	instance, err = tsuru.GetInstanceByName("name")
	c.Assert(err, check.IsNil)
	c.Assert(instance.Apps, check.HasLen, 0)
}

func (s *S) TestServiceUnbindAppWithWizard(c *check.C) {
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
		Name:      "name",
		ScaleUp:   scaleUp,
		ScaleDown: scaleDown,
		Process:   "web",
	}
	err := wizard.New(autoScale)
	c.Assert(err, check.IsNil)
	service := &tsuru.Instance{
		Name: "name",
	}
	err = tsuru.NewInstance(service)
	c.Assert(err, check.IsNil)
	instance, err := tsuru.GetInstanceByName("name")
	c.Assert(err, check.IsNil)
	err = instance.AddApp("", "tsuru-dashboard.192.168.50.4.nip.io")
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	body := `app-host=tsuru-dashboard.192.168.50.4.nip.io`
	request, err := http.NewRequest("DELETE", "/resources/name/bind-app", strings.NewReader(body))
	c.Assert(err, check.IsNil)
	request.Header.Add("Authorization", "token")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	instance, err = tsuru.GetInstanceByName("name")
	c.Assert(err, check.IsNil)
	c.Assert(instance.Apps, check.HasLen, 0)
	_, err = wizard.FindByName(autoScale.Name)
	c.Assert(err, check.NotNil)
}

func (s *S) TestServiceRemove(c *check.C) {
	service := &tsuru.Instance{
		Name: "name",
	}
	err := tsuru.NewInstance(service)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("DELETE", "/resources/name", nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	instance, err := tsuru.GetInstanceByName("name")
	c.Assert(err, check.NotNil)
	c.Assert(instance, check.IsNil)
}

func (s *S) TestServiceInstances(c *check.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"Name":"instance"}]`))
	}))
	defer ts.Close()
	err := os.Setenv("TSURU_HOST", ts.URL)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/service/instance", nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	c.Assert(recorder.HeaderMap["Content-Type"], check.DeepEquals, []string{"application/json"})
	body := recorder.Body.Bytes()
	var instances []tsuru.Instance
	err = json.Unmarshal(body, &instances)
	c.Assert(err, check.IsNil)
	c.Assert(instances, check.HasLen, 1)
	c.Assert(instances[0].Name, check.Equals, "instance")
}

func (s *S) TestServiceInstanceByName(c *check.C) {
	i := &tsuru.Instance{
		Name: "instance",
	}
	err := tsuru.NewInstance(i)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/service/instance/instance", nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	c.Assert(recorder.HeaderMap["Content-Type"], check.DeepEquals, []string{"application/json"})
	body := recorder.Body.Bytes()
	var instance tsuru.Instance
	err = json.Unmarshal(body, &instance)
	c.Assert(err, check.IsNil)
	c.Assert(instance.Name, check.Equals, "instance")
}
