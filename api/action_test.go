// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
	"github.com/tsuru/tsuru-autoscale/action"
	"gopkg.in/check.v1"
)

func server(w http.ResponseWriter, r *http.Request) {
	m := mux.NewRouter()
	Router(m)
	m.ServeHTTP(w, r)
}

func (s *S) TestNewAction(c *check.C) {
	body := `{"name":"new","url":"http://tsuru.io","method":"GET"}`
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/action", strings.NewReader(body))
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusCreated)
}

func (s *S) TestAllActions(c *check.C) {
	err := action.New(&action.Action{URL: "http://tsuru.io", Method: "GET"})
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/action", nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	c.Assert(recorder.HeaderMap["Content-Type"], check.DeepEquals, []string{"application/json"})
	body := recorder.Body.Bytes()
	var a []action.Action
	err = json.Unmarshal(body, &a)
	c.Assert(err, check.IsNil)
	c.Assert(a, check.HasLen, 1)
}

func (s *S) TestRemoveActionNotFound(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("DELETE", "/action", nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusNotFound)
}

func (s *S) TestRemoveAction(c *check.C) {
	a := &action.Action{Name: "namezito", URL: "http://tsuru.io", Method: "GET"}
	err := action.New(a)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("DELETE", fmt.Sprintf("/action/%s", a.Name), nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	_, err = action.FindByName(a.Name)
	c.Assert(err, check.NotNil)
}

func (s *S) TestActionInfo(c *check.C) {
	a := &action.Action{URL: "http://tsuru.io", Method: "GET", Name: "some"}
	err := action.New(a)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", fmt.Sprintf("/action/%s", a.Name), nil)
	request.Header.Add("Authorization", "token")
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	c.Assert(recorder.HeaderMap["Content-Type"], check.DeepEquals, []string{"application/json"})
	body := recorder.Body.Bytes()
	var got action.Action
	err = json.Unmarshal(body, &got)
	c.Assert(err, check.IsNil)
	c.Assert(a.Name, check.Equals, got.Name)
}
