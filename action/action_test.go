// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package action

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

func (s *S) TestNew(c *check.C) {
	url, err := url.Parse("http://tsuru.io")
	c.Assert(err, check.IsNil)
	a, err := New("action", url)
	c.Assert(err, check.IsNil)
	c.Assert(a.Name, check.Equals, "action")
	c.Assert(a.URL, check.Equals, url.String())
}

func (s *S) TestDo(c *check.C) {
	var called bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()
	url, err := url.Parse(ts.URL)
	a, err := New("action", url)
	c.Assert(err, check.IsNil)
	err = a.Do()
	c.Assert(err, check.IsNil)
	c.Assert(called, check.Equals, true)
}
