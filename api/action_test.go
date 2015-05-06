// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/tsuru/tsuru-autoscale/action"
	"gopkg.in/check.v1"
)

func (s *S) TestNewAction(c *check.C) {
	body := `{"name":"new","url":"http://tsuru.io","method":"GET"}`
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/action", strings.NewReader(body))
	c.Assert(err, check.IsNil)
	r := Router()
	r.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusCreated)
}

func (s *S) TestAllActions(c *check.C) {
	err := action.New(&action.Action{URL: "http://tsuru.io", Method: "GET"})
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/action", nil)
	c.Assert(err, check.IsNil)
	r := Router()
	r.ServeHTTP(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusOK)
	c.Assert(recorder.HeaderMap["Content-Type"], check.DeepEquals, []string{"application/json"})
	body := recorder.Body.Bytes()
	var a []action.Action
	err = json.Unmarshal(body, &a)
	c.Assert(err, check.IsNil)
	c.Assert(a, check.HasLen, 1)
}
