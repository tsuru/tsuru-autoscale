// Copyright 2016 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"net/http"
	"net/http/httptest"

	"github.com/tsuru/tsuru-autoscale/action"
	"gopkg.in/check.v1"
)

func (s *S) TestActionRemove(c *check.C) {
	a := &action.Action{
		Name:   "myaction",
		URL:    "http://tsuru.io",
		Method: "GET",
	}
	err := action.New(a)
	c.Assert(err, check.IsNil)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/action/myaction/delete", nil)
	c.Assert(err, check.IsNil)
	server(recorder, request)
	c.Assert(recorder.Code, check.Equals, http.StatusFound)
	_, err = action.FindByName(a.Name)
	c.Assert(err, check.NotNil)
}
