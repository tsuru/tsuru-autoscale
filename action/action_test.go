// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package action

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

func (s *S) TestNew(c *check.C) {
	a := Action{URL: "http://tsuru.io", Method: "GET"}
	err := New(&a)
	c.Assert(err, check.IsNil)
}

func (s *S) TestDo(c *check.C) {
	var called bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()
	a := Action{URL: ts.URL, Method: "GET"}
	err := New(&a)
	c.Assert(err, check.IsNil)
	err = a.Do()
	c.Assert(err, check.IsNil)
	c.Assert(called, check.Equals, true)
}
