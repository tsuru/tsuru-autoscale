// Copyright 2015 tsuru-autoscale authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package action

import (
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
	c.Assert(a.URL, check.Equals, url)
}
