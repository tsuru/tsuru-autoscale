// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"gopkg.in/check.v1"
)

func (s *S) TestServiceAdd(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/resources", nil)
	c.Assert(err, check.IsNil)
	r := Router()
	r.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusCreated)
}

func (s *S) TestServiceBindUnit(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/resources/name/bind", nil)
	c.Assert(err, check.IsNil)
	r := Router()
	r.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusCreated)
}

func (s *S) TestServiceBindApp(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/resources/name/bind-app", nil)
	c.Assert(err, check.IsNil)
	r := Router()
	r.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusCreated)
	var i interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &i)
	c.Assert(err, check.IsNil)
}

func (s *S) TestServiceUnbindUnit(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("DELETE", "/resources/name/bind", nil)
	c.Assert(err, check.IsNil)
	r := Router()
	r.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
}

func (s *S) TestServiceUnbindApp(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("DELETE", "/resources/name/bind-app", nil)
	c.Assert(err, check.IsNil)
	r := Router()
	r.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
}

func (s *S) TestServiceRemove(c *check.C) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("DELETE", "/resources/name", nil)
	c.Assert(err, check.IsNil)
	r := Router()
	r.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
}
